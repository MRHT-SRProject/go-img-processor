#include "img_processor.h"
#include <opencv2/core/cvstd.hpp>
#include <opencv2/core/mat.hpp>
#include <opencv2/imgcodecs/imgcodecs.hpp>
#include <opencv2/imgproc/imgproc.hpp>
#include <opencv2/photo/photo.hpp>
#include <opencv2/video/tracking.hpp>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <vector>

CPixels cvMatToCPixels(cv::Mat);

CMat colorize(CMat img, uint8_t map)
{
    cv::Mat *colorized = new cv::Mat;
    cv::applyColorMap(*((cv::Mat *)img.mat), *colorized, map);

    CMat cmat = CMat{
        .mat = colorized};

    return cmat;
}

CMat stackImages(Image *img, size_t len)
{
    cv::Mat firstImg;
    cv::Mat stacked;
    for (int i = 0; i < len; i++)
    {
        // Transformation matrix
        cv::Mat M = cv::Mat::eye(3, 3, CV_32F);
        cv::Mat pixels = cv::Mat(img[i].y1 - img[i].y0, img[i].x1 - img[i].x0, CV_8UC1, img[i].pixels);
        if (!i)
        {
            firstImg = pixels;
            stacked = pixels;
            continue;
        }
        cv::Mat warped;
        // create the transformation matrix
        const int iterations = 100;         // 300 iterations for alignment
        const double terminationEps = 1e-6; // Threshold in correleation
        // define termination criteria
        cv::TermCriteria criteria(cv::TermCriteria::COUNT + cv::TermCriteria::EPS, iterations, terminationEps);
        cv::findTransformECC(firstImg, pixels, M, cv::MOTION_HOMOGRAPHY, criteria);
        // warp the image according ot the transformation matrix
        cv::warpPerspective(pixels, warped, M, pixels.size(), cv::INTER_LINEAR + cv::WARP_INVERSE_MAP);
        // stack the image
        stacked += warped;
    }

    // cv::applyColorMap(stacked, colorized, cv::COLORMAP_PLASMA);
    CMat cmat = CMat{
        .mat = new cv::Mat(stacked.clone())
    };

    return cmat;
}

CPixels getPixels(CMat mat)
{
    cv::Mat cvmat = *(cv::Mat *)(mat.mat);
    return cvMatToCPixels(cvmat);
}

CMat grayscale(Image img)
{
    cv::Mat pixels = cv::Mat(img.y1 - img.y0, img.x1 - img.x0, CV_8UC4, img.pixels);
    cv::Mat *gray = new cv::Mat();
    cv::cvtColor(pixels, *gray, cv::COLOR_RGBA2GRAY);
    return CMat{.mat = gray};
}

void cMatToImg(CMat mat, const char* filename)
{
    cv::Mat *m = (cv::Mat *)(mat.mat);
    cv::imwrite(filename, *m);
}

CPixels cvMatToCPixels(cv::Mat mat)
{
    return (CPixels){
        .pixels = mat.data,
        .len = mat.total(0)};
}