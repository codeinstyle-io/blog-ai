document.addEventListener("DOMContentLoaded", () => {
    // Create scroll progress bar
    const progressBar = document.createElement("div");
    progressBar.className = "scroll-progress";
    document.body.appendChild(progressBar);

    // Update scroll progress
    window.addEventListener("scroll", () => {
        const windowHeight = window.innerHeight;
        const documentHeight = document.documentElement.scrollHeight -
            windowHeight;
        const scrolled = window.scrollY;
        const width = `${(scrolled / documentHeight) * 100}%`;
        progressBar.style.width = width;

        // Reveal elements on scroll
        document.querySelectorAll(".reveal").forEach((element) => {
            const elementTop = element.getBoundingClientRect().top;
            if (elementTop < windowHeight - 100) {
                element.classList.add("active");
            }
        });
    });

    // Assign random progress values to progress bars
    document.querySelectorAll(".progress-bar").forEach((bar) => {
        const randomValue = Math.floor(Math.random() * 41) + 60; // Random value between 60 and 100
        const progressBarLength = 20; // Length of the progress bar
        const filledLength = Math.round((randomValue / 100) * progressBarLength);
        const emptyLength = progressBarLength - filledLength;
        bar.textContent = `[${'█'.repeat(filledLength)}${' '.repeat(emptyLength)}] ${randomValue}%`;
    });
});

document.addEventListener("DOMContentLoaded", () => {
    const h1 = document.querySelector(".hero h1");
    const text = "Code In Style";
    h1.textContent = "";

    let charIndex = 0;
    function typeText() {
        if (charIndex < text.length) {
            h1.textContent += text.charAt(charIndex);
            charIndex++;
            setTimeout(typeText, 150); // Adjust speed by changing timeout
        }
    }

    typeText();
});
