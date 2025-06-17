// Function to show toast notification
function showToast(message, type = "success") {
  const toast = document.getElementById("toast");
  const toastMessage = document.getElementById("toast-message");
  const icon = toast.querySelector("i");

  // Set message and icon
  toastMessage.textContent = message;
  icon.className =
    type === "success"
      ? "fas fa-check-circle text-code-green mr-2"
      : "fas fa-exclamation-circle text-code-red mr-2";

  // Show toast
  toast.classList.remove("translate-y-full", "opacity-0");
  toast.classList.add("translate-y-0", "opacity-100");

  // Hide toast after 3 seconds
  setTimeout(() => {
    toast.classList.remove("translate-y-0", "opacity-100");
    toast.classList.add("translate-y-full", "opacity-0");
  }, 3000);
}

// Function to load environment variables
async function loadVariables() {
  try {
    const response = await fetch("/api/v1/env");
    const variables = await response.json();

    const container = document.getElementById("variables-container");
    container.innerHTML = "";

    Object.entries(variables).forEach(([key, value]) => {
      addVariableToContainer(key, value);
    });
  } catch (error) {
    console.error("Error loading variables:", error);
    showToast("Error loading environment variables", "error");
  }
}

// Function to add a new variable to the container
function addVariableToContainer(key = "", value = "") {
  const template = document.getElementById("variable-template");
  const container = document.getElementById("variables-container");

  const clone = template.content.cloneNode(true);
  const keyInput = clone.querySelector(".variable-key");
  const valueInput = clone.querySelector(".variable-value");

  keyInput.value = key;
  valueInput.value = value;

  // Add animation class
  const variableItem = clone.querySelector(".variable-item");
  variableItem.classList.add("animate-fade-in");

  container.appendChild(clone);
}

// Function to add a new variable
function addNewVariable() {
  addVariableToContainer();
}

// Function to remove a variable
function removeVariable(button) {
  const item = button.closest(".variable-item");
  item.classList.add("animate-fade-out");

  // Wait for animation to complete before removing
  setTimeout(() => {
    item.remove();
  }, 300);
}

// Function to save variables
async function saveVariables() {
  const variables = {};
  const items = document.querySelectorAll(".variable-item");

  items.forEach((item) => {
    const key = item.querySelector(".variable-key").value.trim();
    const value = item.querySelector(".variable-value").value.trim();

    if (key) {
      variables[key] = value;
    }
  });

  try {
    const response = await fetch("/api/v1/env", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(variables),
    });

    if (response.ok) {
      showToast("Variables saved successfully");
      loadVariables(); // Reload to ensure synchronization
    } else {
      throw new Error("Error saving");
    }
  } catch (error) {
    console.error("Error saving variables:", error);
    showToast("Error saving environment variables", "error");
  }
}

// Add animation styles
const style = document.createElement("style");
style.textContent = `
    @keyframes fadeIn {
        from { opacity: 0; transform: translateY(10px); }
        to { opacity: 1; transform: translateY(0); }
    }
    @keyframes fadeOut {
        from { opacity: 1; transform: translateY(0); }
        to { opacity: 0; transform: translateY(10px); }
    }
    .animate-fade-in {
        animation: fadeIn 0.3s ease-out;
    }
    .animate-fade-out {
        animation: fadeOut 0.3s ease-out;
    }
`;
document.head.appendChild(style);

// Load variables on startup
document.addEventListener("DOMContentLoaded", loadVariables);
