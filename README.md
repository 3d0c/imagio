Image processing web service.
-----------------------------

#### It uses:
  - [OpenCV](http://opencv.org/) for image processing.
  - [GroupCache](https://github.com/golang/groupcache) as a storage backend.
  - [Go http server](http://golang.org/pkg/net/http/) to serve content.

#### Right now it can:  
  - Scale images
  - Crop images
  - Convert formats (jpg,png)
   
#### It's pretty fast.
  - Thanks to OpenCV, it could resize (downscale) up to 140 FullHD images per second. (tested on 3.2Ghz Xeon).
  - Thanks to GroupCache and Go http server, it could serve up to 20k requests per second.

#### It scales well.
  - Thanks to GroupCache and Go http server, it could store unlimited count of images. (Just add a nodes).

Installation.
-------------
### Prerequisite.
+ [Mercurial](http://mercurial.selenic.com/)
+ [Go](http://golang.org/)
+ [Gcc](http://gcc.gnu.org/) — You will need C and C++ compilers.

### 1. Install GroupCache.
```sh
go get github.com/golang/groupcache
```
  
### 2. Building `imagio`.
There are two ways to build it, using the statically linked OpenCV libraries, which are attached to the release package or build it self.  
If You don't have OpenCV and don't plan to use it further, just follow step 2.1.
 
### 2.1 Using statically linked OpenCV libraries from the package.
+ Download the package
  - [OSX version, imagio v.0.1](https://github.com/3d0c/imagio/releases/download/0.1/imagio-0.1-247.osx.tar.gz)
  - [Linux version, imagio v.0.1](https://github.com/3d0c/imagio/releases/download/0.1/imagio-0.1-247.lin.tar.gz)
  
Be sure, that GOROOT and GOPATH variables are set, and run script inside the imagio directory.
```sh
~/imagio# ./build-static.sh
```
That's all. You will get ready to use `imagio` binary. 

### 2.2 Using shared OpenCV libraries
Install OpenCV and ensure that pkgconfig file is available, add it to PKG_CONFIG_PATH if needed.
```sh
# check it
pkg-config --libs opencv
```
```sh
# If You see an error about 'opencv.pc', run the following command
# with corresponding opencv path:
export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/local/opencv-2.4.7/lib/pkgconfig

```  

In couse of this bug [Bug #1925](http://code.opencv.org/issues/1925), You should patch opencv.pc by running following command:
```sh
# copy-paste it
pcPrefix=`grep "prefix=" $PKG_CONFIG_PATH/opencv.pc | grep -v exec | sed 's/prefix=//g'`;pcLibs=`grep "Libs: " $PKG_CONFIG_PATH/opencv.pc`" -L$pcPrefix/lib";sed -i.old 's#libdir=#libdir='"$pcPrefix/lib"'#g' $PKG_CONFIG_PATH/opencv.pc;sed -i.old 's#Libs:.*#'$pcLibs'#g' $PKG_CONFIG_PATH/opencv.pc
```  

Install package by running:
```sh
go get github.com/3d0c/imagio
```

Usage.
------
Example:
```sh
http://localhost:15900/?scale=800x&quality=80&source=farm5.staticflickr.com/4130/5088414872_0856bb93ed_o.jpg
```
As a result You will get a downscaled 800 px width jpeg, saved with 80% quality.  

### Available options:
1. **source**  
  Possible values:
  + `http://some.host.com/1.jpg`
  + `some.host.com/1.jpg` — scheme could be ommitted, http scheme is default)
  + `1.jpg` — host could be ommitted, if it given in config
  + `file://some/path/1.jpg` root option should be defined in config file, default is `/tmp`

2. **scale**  
  Prototype: `([0-9]+x) or (x[0-9]+) or ([0-9]+) or (0.[0-9]+)`  
  E.g.:
  + 800x  width 800px, height will be calculated
  + x600  height 600px, width will be calculated
  + 640   maximum dimension is 640px, e.g. original 1024x768 pixel image will be scaled to 640x480,
          same option applied for 900x1600 image results 360x640
  + 0.5   50% of original dimensions, e.g. 1024x768 = 512x384

3. **crop**  
  Prototype: `crop=x,y,width,height`  
  + `x,y` is a coordinates of top left corner of crop ROI and could be replaced by one of following shortcuts:
    - `left`
    - `bleft`
    - `right`
    - `bright`
    - `center`
  + E.g:
    - &crop=15,20,200,200
    - &crop=center,500,500

4. **quality** Integer value from 0 to 100. (more is better)

5. **format** `jpg` or `png`. Could be omitted if no format conversion needed.

6.  **method** Scaling method. Default is Bicubic  
   Possible values:
    - `1` Nearest-neighbor interpolation
    - `2` Bilinear interpolation
    - `3` Bicubic interpolation
    - `4` Area based interpolation
    - `5` Lanczos resampling

### imagio.conf
If You need to change some default behavior, create an imagio.conf by running:
```
imagio -dumpcfg
```
It will create a default config file in the same directory:
```json
{
    "listen": "127.0.0.1:15900",
    "source": {
        "http": {
            "root": "",
            "default": true
        },
        "file": {
            "root": "",
            "default": false
        }
    },
    
    "defaults": {
        "format": "jpeg",
        "method": 3,
        "quality": 80
    },
    
    "groupcache": {
        "self": "http://127.0.0.1:9100",
        "peers": [],
        "size": "512M"
    }
}
```
It's pretty straightforward. Few comments:
- to use local files, You should setup the `root` option in `file` section
- Groupcache `peers` is an array of strings, e.g. `"peers" : ["host1:9100", "host2:9100"]`
- Groupcache `size` option supports `M` for Megabytes and `G` for Gigabytes
