// Convenience shortcut: Submit the current form when Ctrl-Enter is pressed.
document.addEventListener('keydown', function(e) {
    if (e.keyCode == 13 && e.ctrlKey && e.target.form) {
	e.target.form.submit();
    }
});

function toggleVisibility(id) {
    var e = document.getElementById(id);
    if (e.style.display == "block") {
        e.style.display = "none";
    } else {
        e.style.display = "block";
    }
}
