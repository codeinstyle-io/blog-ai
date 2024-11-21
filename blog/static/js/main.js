(function() {

    function displayProgressBar() {
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
        });
    }

    function displaySkills() {
        // Only execute skills-related code if the container exists
        const skillsContainer = document.getElementById('skills-container');
        if (!skillsContainer) {
            return
        }

        // Fetch skills data from skills.json
        fetch('/skills')
            .then(response => response.json())
            .then(data => {
                data.forEach(section => {
                    // Create section element
                    const sectionElement = document.createElement('section');
                    const sectionTitle = document.createElement('p');
                    sectionTitle.textContent = section.section + ':';
                    sectionElement.appendChild(sectionTitle);

                    // Create skills
                    section.skills.forEach(skill => {
                        const techItem = document.createElement('div');
                        techItem.classList.add('tech-item');

                        const techName = document.createElement('div');
                        techName.classList.add('tech-name');
                        techName.textContent = skill.name;
                        techItem.appendChild(techName);

                        const progressBarContainer = document.createElement('div');
                        progressBarContainer.classList.add('progress-bar-container');

                        const progressBarFill = document.createElement('div');
                        progressBarFill.classList.add('progress-bar-fill');
                        progressBarFill.style.width = '0%';

                        const percentageLabel = document.createElement('span');
                        percentageLabel.classList.add('progress-bar-percentage');
                        percentageLabel.textContent = skill.value + '%';
                        progressBarContainer.appendChild(progressBarFill);
                        progressBarContainer.appendChild(percentageLabel);
                        
                        techItem.appendChild(progressBarContainer);
                        sectionElement.appendChild(techItem);

                        setTimeout(() => {
                            progressBarFill.style.width = skill.value + '%';
                        }, 100);
                    });

                    skillsContainer.appendChild(sectionElement);
                });
            })
            .catch(error => {
                console.error('Error fetching skills data:', error);
            });
    }

    function typeTitle() {
        const h1Span = document.querySelector("h1.type-title span");

        if(!h1Span) {
            return;
        }

        const text = "Code In Style";
        h1Span.textContent = "";

        let charIndex = 0;
        function typeText() {
            if (charIndex < text.length) {
                h1Span.textContent += text.charAt(charIndex);
                charIndex++;
                setTimeout(typeText, 150); // Adjust speed by changing timeout
            }
        }

        typeText();
    }

    function initializeEditor() {
        const editor = document.getElementById('content');
        const preview = document.getElementById('preview-area');
        const editBtn = document.getElementById('edit-mode');
        const previewBtn = document.getElementById('preview-mode');
        
        if (!editor) return;
    
        // Configure marked options
        marked.setOptions({
            gfm: true,
            breaks: true,
            highlight: function(code) {
                return code;
            }
        });
    
        // Live preview
        editor.addEventListener('input', () => {
            preview.innerHTML = marked.parse(editor.value);
        });
    
        // Toggle preview mode
        editBtn.addEventListener('click', () => {
            editor.style.display = 'block';
            preview.style.display = 'none';
            editBtn.classList.add('active');
            previewBtn.classList.remove('active');
        });
    
        previewBtn.addEventListener('click', () => {
            editor.style.display = 'none';
            preview.style.display = 'block';
            editBtn.classList.remove('active');
            previewBtn.classList.add('active');
            preview.innerHTML = marked.parse(editor.value);
        });
    
        // Auto-generate slug from title
        const titleInput = document.getElementById('title');
        const slugInput = document.getElementById('slug');
        
        titleInput.addEventListener('input', () => {
            slugInput.value = titleInput.value
                .toLowerCase()
                .replace(/[^a-z0-9]+/g, '-')
                .replace(/(^-|-$)/g, '');
        });
    }

    displayProgressBar();

    document.addEventListener("DOMContentLoaded", () => {
        displaySkills();
        typeTitle();
        initializeEditor();
    });

}());