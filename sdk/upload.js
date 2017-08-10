/**
 *
 * @param file
 * 文件二进制Blob
 * @constructor
 */
function FileSlicer(file) {
    this.sliceSize = 1024 * 1024 * 5;
    this.slices = Math.ceil(file.size / this.sliceSize);

    this.currentSlice = 0;

    this.getNextSlice = function () {
        var start = this.currentSlice * this.sliceSize
        var end = Math.min((this.currentSlice + 1) * this.sliceSize, file.size);
        ++this.currentSlice;

        return file.slice(start, end);
    }


}

/**
 *
 * @param url
 * ws上传地址 FrozenGo daemon 默认监听 /upload
 * @param file
 * 文件名称
 * @param key
 * 由panel分发的key
 * @param name
 * 文件上传后的名称
 * @param mode
 * 文件上传后的文件属性
 * @constructor
 */
function Uploader(url, file, key, name, mode) {
    var fs = new FileSlicer(file);
    var ws = new WebSocket(url);


    ws.onopen = function () {
        ws.send(key);
        ws.onmessage = function (ms) {
            if (ms.data == "OK") {
                fs.slices--;
                if (fs.slices > 0) ws.send(fs.getNextSlice()); else ws.close();
            } else if (ms.data == "Verified key"){
                ws.send("/" + name + "|" + mode);
            } else if(ms.data == "Ready to upload"){
                ws.send(fs.getNextSlice())
            } else {
                ws.close();
            }
        }
    }

}
