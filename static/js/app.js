/*******************
 *  Element Refs
 *******************/
// Navbar & Buttons
const openAuthBtn = document.getElementById("openAuthBtn");
const authModal = document.getElementById("authModal");
const closeModalBtn = document.getElementById("closeModalBtn");

// Login Form
const loginContainer = document.getElementById("loginContainer");
const loginForm = document.getElementById("loginForm");
const loginEmail = document.getElementById("loginEmail");
const loginPassword = document.getElementById("loginPassword");
const loginSubmit = document.getElementById("loginSubmit");

// Sign Up Form
const signUpContainer = document.getElementById("signUpContainer");
const signUpForm = document.getElementById("signUpForm");
const signUpEmail = document.getElementById("signUpEmail");
const signUpUsername = document.getElementById("signUpUsername");
const signUpPassword = document.getElementById("signUpPassword");
const signUpSubmit = document.getElementById("signUpSubmit");

// Switch links
const showSignUpLink = document.getElementById("showSignUpLink");
const showLoginLink = document.getElementById("showLoginLink");

// Posts
const postsContainer = document.getElementById("postsContainer");

// “+” Floating Button & New Post Modal
const fabAddPost = document.getElementById("fabAddPost");
const newPostModal = document.getElementById("newPostModal");
const closeNewPostModalBtn = document.getElementById("closeNewPostModal");
const newPostForm = document.getElementById("newPostForm");
const postTitleInput = document.getElementById("postTitle");
const postContentInput = document.getElementById("postContent");

/*****************************************
 * 1. On Page Load, Check Session + Load Posts
 *****************************************/
document.addEventListener("DOMContentLoaded", async () => {
  // 1) Check session
  try {
    const res = await fetch("/session-check");
    if (res.ok) {
      const data = await res.json(); // e.g. { loggedIn:true, username:"Alice", profilePic:"..." }
      if (data.loggedIn) {
        // Show logged in nav
        showLoggedInNav(data.username, data.profilePic);
        // Show “+” button
        fabAddPost.classList.remove("hidden");
      }
    }
    // If not OK (401, etc.), user not logged in => do nothing (keep default nav)
  } catch (err) {
    console.error("Session check error:", err);
  }

  // 2) Load posts
  loadPosts();
});

/****************************************
 * 2. Show/Hide Auth Modal (Login & Sign Up)
 ****************************************/
openAuthBtn?.addEventListener("click", () => {
  authModal.classList.remove("hidden");
  loginContainer.classList.remove("hidden");
  signUpContainer.classList.add("hidden");
});

closeModalBtn?.addEventListener("click", () => {
  authModal.classList.add("hidden");
});

// Close modal if click outside of the dialog
window.addEventListener("click", (e) => {
  if (e.target === authModal) {
    authModal.classList.add("hidden");
  }
});

// Switch between login & sign up
showSignUpLink?.addEventListener("click", (e) => {
  e.preventDefault();
  loginContainer.classList.add("hidden");
  signUpContainer.classList.remove("hidden");
});
showLoginLink?.addEventListener("click", (e) => {
  e.preventDefault();
  signUpContainer.classList.add("hidden");
  loginContainer.classList.remove("hidden");
});

/************************************************
 * 3. Enable/Disable Login & Sign Up Buttons
 ************************************************/
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

/********************************************
 * 4. Handle Login Form Submission
 ********************************************/
loginForm.addEventListener("submit", async (event) => {
  event.preventDefault();

  const email = loginEmail.value.trim();
  const password = loginPassword.value.trim();

  // Clear previous error message
  const errorMsg = document.getElementById("loginErrorMsg");
  if (errorMsg) errorMsg.remove();

  try {
    const res = await fetch("/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });

    if (!res.ok) {
      const errorText = await res.text();

      // Display the error message under the inputs
      const errorElement = document.createElement("div");
      errorElement.id = "loginErrorMsg";
      errorElement.className = "error-message"; // Use this class for styling
      errorElement.textContent = errorText;

      loginForm.appendChild(errorElement); // Add it below the form
    } else {
      // Expect JSON like { "message": "Login successful", "username": "Alice", "profilePic": "..." }
      const data = await res.json();
      alert(data.message);

      // Close login modal
      authModal.classList.add("hidden");
      // Show “+” button
      fabAddPost.classList.remove("hidden");

      // Switch nav to logged in
      showLoggedInNav(data.username, data.profilePic);
    }
  } catch (err) {
    console.error("Error while logging in:", err);

    const errorElement = document.createElement("div");
    errorElement.id = "loginErrorMsg";
    errorElement.className = "error-message";
    errorElement.textContent = "Network error occurred";

    loginForm.appendChild(errorElement);
  }

  // Reset form
  loginForm.reset();
  validateLogin();
});

/********************************************
 * 5. Handle Sign Up Form Submission
 ********************************************/
signUpForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const email = signUpEmail.value.trim();
  const username = signUpUsername.value.trim();
  const password = signUpPassword.value.trim();

  // Clear previous error message
  const errorMsg = document.getElementById("signUpErrorMsg");
  if (errorMsg) errorMsg.remove();

  try {
    const res = await fetch("/signup", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, username, password }),
    });

    if (!res.ok) {
      const errorText = await res.text();

      // Display the error message under the inputs
      const errorElement = document.createElement("div");
      errorElement.id = "signUpErrorMsg";
      errorElement.className = "error-message"; // Use this class for styling
      errorElement.textContent = errorText;

      signUpForm.appendChild(errorElement); // Add it below the form
    } else {
      const successText = await res.text();

      // Show welcome popup for first-time signup
      showWelcomePopup(username);

      // Switch to login form
      signUpContainer.classList.add("hidden");
      loginContainer.classList.remove("hidden");
    }
  } catch (error) {
    const errorElement = document.createElement("div");
    errorElement.id = "signUpErrorMsg";
    errorElement.className = "error-message";
    errorElement.textContent = "Network error: " + error.message;

    signUpForm.appendChild(errorElement);
  }

  // Reset sign up form
  signUpForm.reset();
  validateSignUp();
});

// Function to show a welcome popup
function showWelcomePopup(username) {
  const popup = document.createElement("div");
  popup.className = "welcome-popup"; // Style it using CSS
  popup.innerHTML = `
    <div class="popup-content">
      <h2>Welcome, ${username}!</h2>
      <p>Thank you for signing up! We're excited to have you on board.</p>
      <button id="closeWelcomePopup">Got it!</button>
    </div>
  `;
  document.body.appendChild(popup);

  // Close the popup on button click
  const closePopupBtn = document.getElementById("closeWelcomePopup");
  closePopupBtn.addEventListener("click", () => {
    popup.remove();
  });
}

/*************************************
 * 6. Load & Display Posts
 *************************************/
async function loadPosts() {
  try {
    const res = await fetch("/get-post", { method: "GET" });
    if (!res.ok) {
      const err = await res.text();
      console.error("Failed to fetch posts:", err);
      return;
    }
    const posts = await res.json();
    renderPosts(posts);
  } catch (error) {
    console.error("Network error fetching posts:", error);
  }
}

function renderPosts(posts) {
  postsContainer.innerHTML = ""; // clear existing

  if (!posts || posts.length === 0) {
    postsContainer.innerHTML = "<p>No posts yet.</p>";
    return;
  }

  posts.forEach((post) => {
    const postDiv = document.createElement("div");
    postDiv.classList.add("post-card");
    postDiv.innerHTML = `
      <h3>${post.title}</h3>
      <p>${post.content}</p>
      <small>Posted by ${post.username} on ${new Date(post.created_at).toLocaleString()}</small>
    `;
    postsContainer.appendChild(postDiv);
  });
}

/*************************************
 * 7. “+” Button & New Post Modal
 *************************************/
fabAddPost?.addEventListener("click", () => {
  newPostModal.classList.remove("hidden");
});

// Close modal when clicking outside of the dialog
window.addEventListener("click", (e) => {
  if (e.target === newPostModal) {
    newPostModal.classList.add("hidden");
    newPostForm.reset(); // Clear the form inputs when closing
  }
});

closeNewPostModalBtn?.addEventListener("click", () => {
  newPostModal.classList.add("hidden");
  newPostForm.reset(); // Clear the form inputs when closing
});

/*************************************
 * 8. Nav UI: Show Logged-In/Out
 *************************************/
function showLoggedInNav(username, profilePicURL) {
  const navBar = document.querySelector(".navbar");
  // Replace the original nav with a “Logout” + username + plus button
  navBar.innerHTML = `
    <h1 class="logo">My Forum</h1>
    <div class="nav-logged-in" style="display:flex;align-items:center;gap:1rem;">
      ${
        profilePicURL
          ? `<img src="${profilePicURL}" alt="Profile" 
               style="width:32px;height:32px;border-radius:50%;object-fit:cover;">`
          : ""
      }
      <span>${username}</span>
      <!-- The plus button next to Logout -->
      <button id="navAddPostBtn" class="nav-btn" title="Add new post">+</button>
      <button id="logoutBtn" class="nav-btn">Logout</button>
    </div>
  `;

  // Show floating + button also, if you want
  fabAddPost.classList.remove("hidden");

  // Add event for the nav plus button
  document.getElementById("navAddPostBtn").addEventListener("click", () => {
    newPostModal.classList.remove("hidden");
  });

  // Attach logout logic
  document.getElementById("logoutBtn").addEventListener("click", async () => {
    try {
      const res = await fetch("/logout", { method: "POST" });
      if (!res.ok) {
        const errText = await res.text();
        alert(`Logout failed: ${errText}`);
        return;
      }
      const msg = await res.text();
      showLoggedOutNav();
    } catch (error) {
      alert("Logout error: " + error.message);
    }
  });
}

function showLoggedOutNav() {
  const navBar = document.querySelector(".navbar");
  navBar.innerHTML = `
    <h1 class="logo">My Forum</h1>
    <button id="openAuthBtn" class="nav-btn">Log In</button>
  `;

  // Hide the floating + button
  fabAddPost.classList.add("hidden");

  // Re-attach openAuthBtn logic
  const openBtn = document.getElementById("openAuthBtn");
  openBtn.addEventListener("click", () => {
    authModal.classList.remove("hidden");
    loginContainer.classList.remove("hidden");
    signUpContainer.classList.add("hidden");
  });
}