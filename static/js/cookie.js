document.addEventListener("DOMContentLoaded", function () {
  const banner = document.getElementById("cookie-banner");
  const acceptButton = document.getElementById("accept-cookies");
  if (!getCookie("cookies-accepted")) {
    banner.style.display = "block";
  }

  acceptButton.addEventListener("click", function () {
    setCookie("cookies-accepted", "true", 365);
    banner.style.display = "none";
  });
});

function setCookie(name, value, days) {
  let expires = "";
  if (days) {
    const date = new Date();
    data.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
    expires = "; expires=" + date.toUTCString();
  }
  document.cookie = name + "=" + (value || "") + expires + "; path=/";
}

function getCookie(name) {
  const nameEq = name + "=";
  const ca = document.cookie.split(";");
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) == " ") c = c.substring(1, c.length);
    if (c.indexOf(nameEq) == 0) return c.substring(nameEq.length, c.length);
  }
  return null;
}
