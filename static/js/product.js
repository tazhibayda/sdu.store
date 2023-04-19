function image(x){
    document.getElementById("main_image_of_product").src = x.src;
}
function color(c){
    document.getElementById("color_title").innerHTML = "Color: " + c;
}
function size(x){
    document.getElementById("size_title").innerHTML = "Size: " + x.innerHTML;
}