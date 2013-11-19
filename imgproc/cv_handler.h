#ifndef __cv_scale_h__
#define __cv_scale_h__

#include <highgui.h>
#include <cv.h>
#include <stdio.h>
#include <strings.h>
#include <stdlib.h>
#include <stdint.h>

typedef struct {
    int64_t width;
    int64_t height;
} PixelDim;

typedef struct {
    int64_t X;
    int64_t Y;
    int64_t width;
    int64_t height;
} ROI;

typedef struct {
    unsigned char *data;
    unsigned int length;    
} Blob;

Blob *cv_handler(Blob *in, PixelDim *zoom, int quality, int method, const char *format, ROI *roi);

#endif
