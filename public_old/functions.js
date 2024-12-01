function account_system() {
    // does not work because jwt_token is HttpOnly
    if (document.cookie("jwt_token") != "") {
        document.querySelector('#logout').style.display = 'none';
        document.querySelector('.account_management').style.display = 'flex';
    } 
}

document.addEventListener("DOMContentLoaded", function() {
    var assignments_nav = document.getElementById("assignments_nav");
    if (assignments_nav) {
        assignments_nav.remove();
    }
});