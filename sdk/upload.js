function FileSlicer(file) {
    this.sliceSize = 1024 * 1024;
    this.slices = Math.ceil(file.size / this.sliceSize);

    this.currentSlice = 0;

    this.getNextSlice = function (){
      var start = (this.currentSlice + 1) * this.sliceSize
      var end = Math.min(this.currentSlice * this.sliceSize,file.size);
      ++this.currentSlice;

      return file.slice(start,end);
    }


}

function Uploader(url,file) {
    var fs = new FileSlicer(file);
    var ws = new WebSocket(url);

    ws.onopen = function () {
        ws.send(fs.getNextSlice());
    }
    ws.onmessage = function (ms) {
        if(ms.data=="OK"){
            fs.slices--
            if(fs.slices > 0) socket.send(fs.getNextSlice())
        } else {

        }
    }
}