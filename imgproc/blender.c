#include "cv_handler.h"

/*
    Alpha blending with a mask, using OpenCV, is pretty simple:

    cvThreshold(mask, mask, 180, 255, 1);
    cvNot(mask, mask);
    cvAdd(result, fg, result, mask);

    but result is ugly,— i couldn't get smooth image from curved masks.

    Here is a bit ported solution of combining BGR background image with transparent foreground
    originally written by Michael Jepson. It's about 5-9ms slower, than native implementation, but
    works well.
*/

CvMat *overlayImage(CvMat *bgMat, IplImage *bgImg, const CvMat *fgMat, const IplImage *fgImg, CvMat *maskMat, CvRect *roi) {    
    CvMat *resultMat = cvCloneMat(bgMat); 

    int loc_x = (roi->x > 0) ? roi->x : 0;
    int loc_y = (roi->y > 0) ? roi->y : 0;

    int y, x, c;
    
    for(y = loc_y; y < bgMat->rows; y++) {
        int fY = y - roi->y;

        if(fY >= fgMat->rows) {
            break;
        }

        for(x = loc_x; x < bgMat->cols; x++) {
            int fX = x - roi->x;

            if(fX >= fgMat->cols) {
                break;
            }

            double opacity;

            if (maskMat) {
                opacity = ((double)maskMat->data.ptr[fY * maskMat->step + fX]) / 255.;
            } else {
                opacity = ((double)fgMat->data.ptr[fY * fgMat->step + fX * fgImg->nChannels + 3]) / 255.;
            }

            for(c = 0; opacity > 0 && c < bgImg->nChannels; c++) {
                unsigned char foregroundPx = fgMat->data.ptr[fY * fgMat->step + fX * fgImg->nChannels + c];
                unsigned char backgroundPx = bgMat->data.ptr[y * bgMat->step + x * bgImg->nChannels + c];
                resultMat->data.ptr[y*resultMat->step + bgImg->nChannels*x + c] = backgroundPx * (1.-opacity) + foregroundPx * opacity;
            }
        }
    }
    
    return resultMat;
}

IplImage *addWeighted(IplImage *base, IplImage *fg, float alpha, CvRect *roi) {
    IplImage *result = cvCloneImage(base);

    cvSetImageROI(result, cvRect(roi->x, roi->y, fg->width, fg->height));
    cvAddWeighted(result, 1.0, fg, alpha, 0.0, result);
    cvResetImageROI(result);
    
    return result;
}

Blob *blender(const Blob *base, const Blob *foreground, const Blob *mask, const int quality, const char *format, const float alpha, CvRect *roi) {
    int p[3] = {CV_IMWRITE_JPEG_QUALITY, quality, 0};

    cvUseOptimized(1);

    // init base
    CvMat *baseBuf = cvCreateMat(1, base->length, CV_8UC1);
    baseBuf->data.ptr = base->data;

    CvMat *baseMat = cvDecodeImageM(baseBuf, CV_LOAD_IMAGE_UNCHANGED);
    IplImage *baseImg = cvDecodeImage(baseBuf, CV_LOAD_IMAGE_UNCHANGED);
    cvReleaseMat(&baseBuf);

    // init foreground
    CvMat *fgBuf = cvCreateMat(1, foreground->length, CV_8UC1);
    fgBuf->data.ptr = foreground->data;

    CvMat *fgMat = cvDecodeImageM(fgBuf, CV_LOAD_IMAGE_UNCHANGED);
    IplImage *fgImg = cvDecodeImage(fgBuf, CV_LOAD_IMAGE_UNCHANGED);
    cvReleaseMat(&fgBuf);

    // init mask
    CvMat *maskMat, *maskBuf = NULL;
    IplImage *maskImg = NULL;
    
    if(mask) {
        maskBuf = cvCreateMat(1, mask->length, CV_8UC1);
        if(maskBuf) {
            maskBuf->data.ptr = mask->data;    
            maskMat = cvDecodeImageM(maskBuf, CV_LOAD_IMAGE_GRAYSCALE);
            maskImg = cvDecodeImage(maskBuf, CV_LOAD_IMAGE_GRAYSCALE);
            cvReleaseMat(&maskBuf);
        }        
    }
    
    CvMat *result;

    if(fgImg->nChannels <= 3 && !mask) {
        IplImage *tmp = addWeighted(baseImg, fgImg, alpha, roi);    
        result = cvEncodeImage(format, tmp, p);
        cvReleaseImage(&tmp);
    } else {
        CvMat *tmp = overlayImage(baseMat, baseImg, fgMat, fgImg, maskMat, roi);    
        result = cvEncodeImage(format, tmp, p);
        cvReleaseMat(&tmp);
    }

    if (mask && maskMat) {
        cvReleaseMat(&maskMat);
        cvReleaseImage(&maskImg);
    }

    cvReleaseMat(&baseMat);
    cvReleaseMat(&fgMat);
    cvReleaseImage(&baseImg);
    cvReleaseImage(&fgImg);

    Blob *out = malloc(sizeof(Blob));
    out->data = malloc(result->step);
    out->length = result->step;

    memcpy(out->data, result->data.ptr, result->step);
    cvReleaseMat(&result);

    return (out);
}
