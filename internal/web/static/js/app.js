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

// Function to load environments and their variables
async function loadEnvironments() {
  try {
    const response = await fetch("/api/v1/env");
    const environments = await response.json();

    const container = document.getElementById("environments-container");
    container.innerHTML = "";

    environments?.data?.forEach((env) => {
      addEnvironmentToContainer(env);
    });
  } catch (error) {
    console.error("Error loading environments:", error);
    showToast("Error loading environments", "error");
  }
}

// Function to add a new environment to the container
function addEnvironmentToContainer(env = null) {
  const template = document.getElementById("environment-template");
  const container = document.getElementById("environments-container");

  const clone = template.content.cloneNode(true);
  const nameInput = clone.querySelector(".environment-name");
  const variablesContainer = clone.querySelector(".variables-container");

  if (env) {
    nameInput.value = env.name;
    nameInput.dataset.id = env.id;

    // Add existing variables
    env.values.forEach((value) => {
      addVariableToContainer(variablesContainer, value);
    });
  }

  // Add animation class
  const environmentItem = clone.querySelector(".environment-item");
  environmentItem.classList.add("animate-fade-in");

  container.appendChild(clone);
}

// Function to add a new environment
function addNewEnvironment() {
  addEnvironmentToContainer();
}

// Function to remove an environment
function removeEnvironment(button) {
  const item = button.closest(".environment-item");
  item.classList.add("animate-fade-out");
  const nameInput = item.querySelector(".environment-name");
  // Wait for animation to complete before removing
  setTimeout(async () => {
    item.remove();
    try {
      const environmentId = nameInput.dataset.id;
      const url = `/api/v1/env/${environmentId}`;

      const response = await fetch(url, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        showToast("Variable deleted successfully");
        loadEnvironments(); // Reload to ensure synchronization
      } else {
        throw new Error("Error deleting variable");
      }
    } catch (error) {
      console.error("Error saving environment:", error);
      showToast("Error saving environment", "error");
    }
  }, 300);
}

// Function to add a new variable to an environment
function addVariableToContainer(container, value = null) {
  const template = document.getElementById("variable-template");
  const clone = template.content.cloneNode(true);
  const keyInput = clone.querySelector(".variable-key");
  const valueInput = clone.querySelector(".variable-value");
  const removeButton = clone.querySelector(".remove-button");

  if (value) {
    keyInput.value = value.key;
    valueInput.value = value.value;
    keyInput.dataset.id = value.id;
    removeButton.dataset.id = value.id;
  }

  // Add animation class
  const variableItem = clone.querySelector(".variable-item");
  variableItem.classList.add("animate-fade-in");

  container.appendChild(clone);
}

// Function to add a new variable
function addNewVariable(button) {
  const environmentItem = button.closest(".environment-item");
  const variablesContainer = environmentItem.querySelector(
    ".variables-container"
  );
  addVariableToContainer(variablesContainer);
}

// Function to remove a variable
function removeVariable(button) {
  const item = button.closest(".variable-item");
  item.classList.add("animate-fade-out");
  // get the environment id
  const environmentItem = button.closest(".environment-item");
  const nameInput = environmentItem.querySelector(".environment-name");
  // get the variable key
  const valueInput = item.querySelector(".remove-button ");

  // Wait for animation to complete before removing
  setTimeout(async () => {
    item.remove();
    try {
      const environmentId = nameInput.dataset.id;
      const valueId = valueInput.dataset.id;
      const url = `/api/v1/env/${environmentId}/value/${valueId}`;

      const response = await fetch(url, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        showToast("Variable deleted successfully");
        loadEnvironments(); // Reload to ensure synchronization
      } else {
        throw new Error("Error deleting variable");
      }
    } catch (error) {
      console.error("Error saving environment:", error);
      showToast("Error saving environment", "error");
    }
  }, 300);
}

// Function to save an environment and its variables
async function saveEnvironment(button) {
  const environmentItem = button.closest(".environment-item");
  const nameInput = environmentItem.querySelector(".environment-name");
  const variablesContainer = environmentItem.querySelector(
    ".variables-container"
  );
  const variables = [];

  variablesContainer.querySelectorAll(".variable-item").forEach((item) => {
    const key = item.querySelector(".variable-key").value.trim();
    const value = item.querySelector(".variable-value").value.trim();

    if (key) {
      variables.push({
        key: key,
        value: value,
      });
    }
  });

  try {
    const environmentId = nameInput.dataset.id;
    const method = environmentId ? "PUT" : "POST";
    const url = environmentId ? `/api/v1/env/${environmentId}` : "/api/v1/env";

    const response = await fetch(url, {
      method: method,
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name: nameInput.value.trim(),
        values: variables,
      }),
    });

    if (response.ok) {
      showToast("Environment saved successfully");
      loadEnvironments(); // Reload to ensure synchronization
    } else {
      throw new Error("Error saving environment");
    }
  } catch (error) {
    console.error("Error saving environment:", error);
    showToast("Error saving environment", "error");
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

// Function to toggle password visibility
function togglePasswordVisibility(button) {
  const input = button.parentElement.querySelector(".variable-value");
  const icon = button.querySelector("i");

  if (input.type === "password") {
    input.type = "text";
    icon.className = "fas fa-eye-slash";
  } else {
    input.type = "password";
    icon.className = "fas fa-eye";
  }
}

// Load environments on startup
document.addEventListener("DOMContentLoaded", loadEnvironments);
