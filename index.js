setInterval(
    function(){
        fetch('http://localhost:3000/metrics')
        .then((resp)=>{
            return resp.json()
        }).then((data)=>{
            console.log(data)
        })
    },1000)

