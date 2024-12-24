// Get references
const openAuthBtn = document.getElementById("openAuthBtn");
const authModal = document.getElementById("authModal");
const closeModalBtn = document.getElementById("closeModalBtn");

// Containers
const loginContainer = document.getElementById("loginContainer");
const signUpContainer = document.getElementById("signUpContainer");

// Switch links
const showSignUpLink = document.getElementById("showSignUpLink");
const showLoginLink = document.getElementById("showLoginLink");

// Login form & fields (not fully implemented in DB for this demo)
const loginForm = document.getElementById("loginForm");
const loginEmail = document.getElementById("loginEmail");
const loginPassword = document.getElementById("loginPassword");
const loginSubmit = document.getElementById("loginSubmit");

// Sign Up form & fields
const signUpForm = document.getElementById("signUpForm");
const signUpEmail = document.getElementById("signUpEmail");
const signUpUsername = document.getElementById("signUpUsername");
const signUpPassword = document.getElementById("signUpPassword");
const signUpSubmit = document.getElementById("signUpSubmit");

// 1. Open/Close Modal
openAuthBtn.addEventListener("click", () => {
  authModal.classList.remove("hidden");
  // Default to login view
  loginContainer.classList.remove("hidden");
  signUpContainer.classList.add("hidden");
});

closeModalBtn.addEventListener("click", () => {
  authModal.classList.add("hidden");
});

// Close if clicking outside modal-dialog
window.addEventListener("click", (e) => {
  if (e.target === authModal) {
    authModal.classList.add("hidden");
  }
});

// 2. Switch between Login and Sign Up
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

// 3. Enable/Disable Submit Buttons
// -- Login (demo only)
function validateLogin() {
  const valid = loginEmail.value.trim() && loginPassword.value.trim();
  loginSubmit.disabled = !valid;
  loginSubmit.classList.toggle("disabled", !valid);
}
loginEmail.addEventListener("input", validateLogin);
loginPassword.addEventListener("input", validateLogin);

// -- Sign Up
function validateSignUp() {
  const valid = signUpEmail.value.trim() &&
                signUpUsername.value.trim() &&
                signUpPassword.value.trim();
  signUpSubmit.disabled = !valid;
  signUpSubmit.classList.toggle("disabled", !valid);
}
signUpEmail.addEventListener("input", validateSignUp);
signUpUsername.addEventListener("input", validateSignUp);
signUpPassword.addEventListener("input", validateSignUp);

// 4. Handle Form Submissions
// -- Login (not implemented in DB)
loginForm.addEventListener("submit", (e) => {
  e.preventDefault();
  alert(`Logging in with email/username: ${loginEmail.value}`);
  // Reset form
  loginForm.reset();
  validateLogin();
  authModal.classList.add("hidden");
});

// -- Sign Up (CREATE user in DB)
signUpForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const email = signUpEmail.value.trim();
  const username = signUpUsername.value.trim();
  const password = signUpPassword.value.trim();

  try {
    const res = await fetch("/signup", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, username, password })
    });

    if (!res.ok) {
      const errorText = await res.text();
      alert(`Sign Up Error: ${errorText}`);
    } else {
      const successText = await res.text();
      alert(successText);
      // Possibly auto-close the modal or switch to login
      authModal.classList.add("hidden");
    }
  } catch (error) {
    alert("Network error: " + error.message);
  }

  // Reset
  signUpForm.reset();
  validateSignUp();
});
