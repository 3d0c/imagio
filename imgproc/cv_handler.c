#include "cv_handler.h"

Blob *cv_handler(Blob *in, PixelDim *zoom, int quality, int method, const char *format, ROI *roi) {
	if (!in) {
		fprintf(stderr, "cv_handler.c: Wrong call. 'in' is NULL\n");
		return NULL;
	}

	if (!zoom && !roi) {
		fprintf(stderr, "cv_handler.c: Wrong call. 'zomm' and 'roi' are NULL\n");
		return NULL;
	}

	IplImage *srcImg, *resultImg;
	int p[3] = {CV_IMWRITE_JPEG_QUALITY, quality, 0};

	cvUseOptimized(1);
	
	CvMat *buf = cvCreateMat(1, in->length, CV_8UC1);
	buf->data.ptr = in->data;

	srcImg = cvDecodeImage(buf, CV_LOAD_IMAGE_COLOR);
	cvReleaseMat(&buf);

	if(!srcImg) {
		fprintf(stderr, "cv_handler.c: cvDecodeImage() error.\n");
		return NULL;
	}

	if(roi) {
		cvSetImageROI(srcImg, cvRect(roi->X, roi->Y, roi->width, roi->height));
	}

	int width = (zoom) ? zoom->width : roi->width;
	int height = (zoom) ? zoom->height : roi->height;

	if(!(resultImg = cvCreateImage(cvSize(width, height), srcImg->depth, srcImg->nChannels))) {
		fprintf(stderr, "cv_handler.c: cvCreateImage error.\n");
		cvReleaseImage(&srcImg);
		return NULL;
	}

	if(!zoom && roi) {
		cvCopyImage(srcImg, resultImg);
	} else {
		cvResize(srcImg, resultImg, method);
	}

	CvMat *result = cvEncodeImage(format, resultImg, p);

	Blob *out = malloc(sizeof(Blob));
	out->data = malloc(result->step);
	out->length = result->step;

	memcpy(out->data, result->data.ptr, result->step);
	
	cvReleaseMat(&result);
	cvReleaseImage(&srcImg);
	cvReleaseImage(&resultImg);
	
	return out;
}
