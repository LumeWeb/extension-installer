// Get input + focus
const installBtn = document.getElementById('install');
installBtn.addEventListener('click', function () {
    installBtn.style.display = 'none';
    window.runtime.EventsEmit("install");
});