var buyButton = document.getElementById("search_button")
var color = ""
var size = ""

buyButton.addEventListener("click", event =>{
    event.preventDefault()
    $.ajax({
        type: 'POST',
        url: '/product/purchase',
        dataType: 'json',
        data:{
            'color': color,
            'size': size,
            'product_id': product_id,
        },
        success: function(response){
            document.getElementById("error").innerHTML = response
        },
        error: function(error){
            document.getElementById("error").innerHTML = error
        }
    })
    var d = document.getElementById("error")
})

function image(x){
    document.getElementById("main_image_of_product").src = x.src;
}
function colorx(c){
    document.getElementById("color_title").innerHTML = "Color: " + c;
    color = c;
}
function sizex(x){
    document.getElementById("size_title").innerHTML = "Size: " + x;
    size = x;
}