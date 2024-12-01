window.addEventListener("click", function(event) {
    // displaying dropdown for profile
    const profileDropdown = document.getElementById('profile-dropdown')
    if (event.target.id == 'profile-btn-text') {
        profileDropdown.style.display = profileDropdown.style.display === 'block' ? 'none' : 'block';
    }
    else {
        profileDropdown.style.display = 'none';
    }

    // adjusting profile button depending on if dropdown is open
    const profileBtn = document.getElementById('profile-btn-text')
    if (profileDropdown.style.display == 'none') {
        profileBtn.textContent = profileBtn.textContent.slice(0, -1) + '⮟';
    }
    else {
        profileBtn.textContent = profileBtn.textContent.slice(0, -1) + '⮝';
    }
});