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

void stackImages(Image *img, size_t len)
{
    cv:: Mat firstImg;
    cv::Mat stacked;
    cv::Mat colorized;
    for (int i = 0; i < len; i++)
    {
        // Transformation matrix
        cv::Mat M = cv::Mat::eye(3, 3, CV_32F);
        // Convert pixels (4 channel RGBA) to Mat
        cv::Mat pixels = cv::Mat(img[i].y1 - img[i].y0, img[i].x1 - img[i].x0, CV_8UC4, img[i].pixels);
        CMat cgray = grayscale(*img);
        cv::Mat gray = cv::Mat(*(cv::Mat*)(cgray.mat ));
        // cleanup
        delete (cv::Mat*)(cgray.mat);
        if(!i) {
            firstImg = gray;
            stacked = gray;
            continue;
        }
        cv::Mat warped;
        // create size struct
        cv::Size size;
        size.width = img[i].x1 - img[i].x0;
        size.height = img[i].y1 - img[i].y0;
        // convert color to bw
        // create the transformation matrix
        const int iterations = 100; // 300 iterations for alignment
        const double terminationEps = 1e-6; // Threshold in correleation
        // define termination criteria
        cv::TermCriteria criteria (cv::TermCriteria::COUNT+cv::TermCriteria::EPS, iterations, terminationEps);
        cv::findTransformECC(firstImg, gray, M, cv::MOTION_HOMOGRAPHY, criteria);
        // warp the image according ot the transformation matrix
        cv::warpPerspective(gray, warped, M, gray.size(), cv::INTER_LINEAR + cv::WARP_INVERSE_MAP);
        // stack the image
        stacked += warped;
    }

    cv::applyColorMap(stacked, colorized, cv::COLORMAP_PLASMA);
    // write the image
    cv::imwrite("stacked.jpg", colorized);


}

CPixels getPixels(CMat mat) {
    cv::Mat cvmat = *(cv::Mat*)(mat.mat);
    return (CPixels){
        .pixels=cvmat.data,
        .len=cvmat.total(0)
    };
} 

CMat grayscale(Image img) {
        cv::Mat pixels = cv::Mat(img.y1 - img.y0, img.x1 - img.x0, CV_8UC4, img.pixels);
        cv::Mat *gray = new cv::Mat();
        cv::cvtColor(pixels, *gray, cv::COLOR_RGBA2GRAY);
        return CMat{.mat = gray};
}