function account_system() {
    if (/* logged in user logged out */) {
        document.querySelector('#logout').style.display = 'none';
        document.querySelector('.account_management').style.display = 'flex';
    } else if (/* user logged in */) {
        document.querySelector('#logout').style.display = 'block';
        document.querySelector('.account_management').style.display = 'none';
    }
}