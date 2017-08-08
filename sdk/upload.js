function FileSlicer(file) {
    this.sliceSize = 1024 * 1024;
    this.slices = Math.ceil(file.size / this.sliceSize);

    this.currentSlice = 0;

    this.getNextSlice = function (){
      var start = this.currentSlice * this.sliceSize
      var end = Math.min((this.currentSlice + 1) * this.sliceSize,file.size);
      ++this.currentSlice;

      return file.slice(start,end);
    }


}

function Uploader(url,file) {
    var fs = new FileSlicer(file);
    var ws = new WebSocket(url)


    ws.onopen = function () {
        ws.send(fs.getNextSlice());
    }
    ws.onmessage = function (ms) {
        if(ms.data=="OK"){
            console.log("get.")
            console.log(ms.data)
            console.log(ms.data == "OK")
            fs.slices--
            if(fs.slices > 0) ws.send(fs.getNextSlice())
        } else {

        }
    }
}
