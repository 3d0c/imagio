#include "cv_handler.h"

Blob *resizer(Blob *in, PixelDim *zoom, int quality, int method, const char *format, CvRect *roi) {
	if (!in) {
		fprintf(stderr, "resizer.c: Wrong call. 'in' is NULL\n");
		return NULL;
	}

	if (!zoom && !roi) {
		fprintf(stderr, "resizer.c: Wrong call. 'zomm' and 'roi' are NULL\n");
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
		fprintf(stderr, "resizer.c: cvDecodeImage() error.\n");
		return NULL;
	}

	if(roi) {
		cvSetImageROI(srcImg, *roi);
	}

	int width = (zoom) ? zoom->width : roi->width;
	int height = (zoom) ? zoom->height : roi->height;

	if(!(resultImg = cvCreateImage(cvSize(width, height), srcImg->depth, srcImg->nChannels))) {
		fprintf(stderr, "resizer.c: cvCreateImage error.\n");
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
