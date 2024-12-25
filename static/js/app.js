// -------------------- Element References --------------------
const openAuthBtn = document.getElementById("openAuthBtn");
const authModal = document.getElementById("authModal");
const closeModalBtn = document.getElementById("closeModalBtn");

// Login container references
const loginContainer = document.getElementById("loginContainer");
const loginForm = document.getElementById("loginForm");
const loginEmail = document.getElementById("loginEmail");
const loginPassword = document.getElementById("loginPassword");
const loginSubmit = document.getElementById("loginSubmit");

// Signup container references
const signUpContainer = document.getElementById("signUpContainer");
const signUpForm = document.getElementById("signUpForm");
const signUpEmail = document.getElementById("signUpEmail");
const signUpUsername = document.getElementById("signUpUsername");
const signUpPassword = document.getElementById("signUpPassword");
const signUpSubmit = document.getElementById("signUpSubmit");

// Switch links
const showSignUpLink = document.getElementById("showSignUpLink");
const showLoginLink = document.getElementById("showLoginLink");

// 1. Open/Close Modal
openAuthBtn.addEventListener("click", () => {
  authModal.classList.remove("hidden");
  loginContainer.classList.remove("hidden");
  signUpContainer.classList.add("hidden");
});

closeModalBtn.addEventListener("click", () => {
  authModal.classList.add("hidden");
});

// Close modal if click outside of the dialog
window.addEventListener("click", (e) => {
  if (e.target === authModal) {
    authModal.classList.add("hidden");
  }
});

// 2. Switch between Login & Sign Up
showSignUpLink.addEventListener("click", (e) => {
  e.preventDefault();
  loginContainer.classList.add("hidden");
  signUpContainer.classList.remove("hidden");
});

showLoginLink.addEventListener("click", (e) => {
  e.preventDefault();
  signUpContainer.classList.add("hidden");
  loginContainer.classList.remove("hidden");
});

// 3. Enable/Disable Login & Sign Up Buttons
function validateLogin() {
  const valid = loginEmail.value.trim() && loginPassword.value.trim();
  loginSubmit.disabled = !valid;
  loginSubmit.classList.toggle("disabled", !valid);
}
loginEmail.addEventListener("input", validateLogin);
loginPassword.addEventListener("input", validateLogin);

function validateSignUp() {
  const valid =
    signUpEmail.value.trim() &&
    signUpUsername.value.trim() &&
    signUpPassword.value.trim();
  signUpSubmit.disabled = !valid;
  signUpSubmit.classList.toggle("disabled", !valid);
}
signUpEmail.addEventListener("input", validateSignUp);
signUpUsername.addEventListener("input", validateSignUp);
signUpPassword.addEventListener("input", validateSignUp);

// 4. Handle Login Form Submission
loginForm.addEventListener("submit", async (event) => {
  event.preventDefault();

  const email = loginEmail.value.trim();
  const password = loginPassword.value.trim();

  try {
    const res = await fetch("/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });

    if (!res.ok) {
      const errorText = await res.text();
      alert(`Login failed: ${errorText}`);
    } else {
      // Expecting JSON like { "message": "Login successful", "username": "Alice", "profilePic": "https://..." }
      const data = await res.json();
      alert(data.message);

      // Close the login modal
      authModal.classList.add("hidden");

      // Show the profile pic & logout in the navbar
      // data.username and data.profilePic come from the server
      showLoggedInNav(data.username, data.profilePic);
    }
  } catch (err) {
    console.error("Error while logging in:", err);
    alert("Network error occurred");
  }

  // Reset form
  loginForm.reset();
  validateLogin();
});

// 5. Handle Sign Up Form Submission
signUpForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const email = signUpEmail.value.trim();
  const username = signUpUsername.value.trim();
  const password = signUpPassword.value.trim();

  try {
    const res = await fetch("/signup", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, username, password }),
    });

    if (!res.ok) {
      const errorText = await res.text();
      alert(`Sign Up Error: ${errorText}`);
    } else {
      const successText = await res.text();
      alert(`Sign Up Success: ${successText}`);

      // Switch to login form if desired
      signUpContainer.classList.add("hidden");
      loginContainer.classList.remove("hidden");
    }
  } catch (error) {
    alert("Network error: " + error.message);
  }

  // Reset sign up form
  signUpForm.reset();
  validateSignUp();
});

// -------------------- Show Logged-In Nav (Profile Pic + Logout) --------------------
function showLoggedInNav(username, profilePicURL) {
  // 1. Remove the old "Log In" button
  const navBar = document.querySelector(".navbar");
  
  // 2. Replace or modify it to show a profile pic + Logout
  // Example approach:
  navBar.innerHTML = `
    <h1 class="logo">My Forum</h1>
    <div class="nav-logged-in">
      
      <span style="margin-left: 8px;">${username}</span>
      <button id="logoutBtn" class="nav-btn" style="margin-left: 16px;">Logout</button>
    </div>
  `;

  // // 3. Add logout event
  // document.getElementById("logoutBtn").addEventListener("click", async () => {
  //   try {
  //     const res = await fetch("/logout", { method: "POST" });
  //     if (!res.ok) {
  //       const errText = await res.text();
  //       alert(`Logout failed: ${errText}`);
  //     } else {
  //       // On success, revert to logged-out nav
  //       const msg = await res.text();
  //       alert(msg);
  //       showLoggedOutNav();
  //     }
  //   } catch (error) {
  //     alert("Logout error: " + error.message);
  //   }
  // });
}

// -------------------- Show Logged-Out Nav (Log In Button) --------------------
function showLoggedOutNav() {
  const navBar = document.querySelector(".navbar");
  navBar.innerHTML = `
    <h1 class="logo">My Forum</h1>
    <button id="openAuthBtn" class="nav-btn">Log In</button>
  `;

  // Re-attach the event for opening the modal
  const openBtn = document.getElementById("openAuthBtn");
  openBtn.addEventListener("click", () => {
    authModal.classList.remove("hidden");
    loginContainer.classList.remove("hidden");
    signUpContainer.classList.add("hidden");
  });
}

// On page load, check if the user is already logged in
document.addEventListener("DOMContentLoaded", async () => {
  try {
    const res = await fetch("/session-check", { method: "GET" });

    if (res.ok) {
      // If 200, the user is logged in
      const data = await res.json(); 
      if (data.loggedIn) {
        // data.username might be returned from the server
        showLoggedInNav(data.username, data.profilePic); 
      } else {
        // Not logged in, do nothing or show default nav
      }
    } else {
      // If 401 or any error, user not logged in
      // Optionally do nothing (show default nav with login)
    }
  } catch (err) {
    console.error("Session check error:", err);
    // If there's a network error, fallback to showing default nav
  }
});
