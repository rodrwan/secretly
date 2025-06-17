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
    alert("Error loading environment variables");
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

  container.appendChild(clone);
}

// Function to add a new variable
function addNewVariable() {
  addVariableToContainer();
}

// Function to remove a variable
function removeVariable(button) {
  button.closest(".variable-item").remove();
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
      alert("Variables saved successfully");
      loadVariables(); // Reload to ensure synchronization
    } else {
      throw new Error("Error saving");
    }
  } catch (error) {
    console.error("Error saving variables:", error);
    alert("Error saving environment variables");
  }
}

// Load variables on startup
document.addEventListener("DOMContentLoaded", loadVariables);
