function account_system() {
    if (document.cookie("jwt_token") != "") {
        document.querySelector('#logout').style.display = 'none';
        document.querySelector('.account_management').style.display = 'flex';
    } else {
        document.querySelector('#logout').style.display = 'block';
        document.querySelector('.account_management').style.display = 'none';
    }
}