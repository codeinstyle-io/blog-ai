function deleteTag(id) {
    fetch(`/admin/tags/${id}`, {
        method: 'DELETE',
    }).then((response) => {
        if (response.ok) {
            window.location.href = '/admin/tags';
        }
    });
}

function deletePost(id) {
    fetch(`/admin/posts/${id}`, {
        method: 'DELETE',
    }).then((response) => {
        if (response.ok) {
            window.location.href = '/admin/posts';
        }
    });
}

function deletePage(id) {
    fetch(`/admin/pages/${id}`, {
        method: 'DELETE',
    }).then(response => {
        if (response.ok) {
            window.location.href = '/admin/pages';
        }
    });
}

function deleteMenuItem(id) {
    fetch(`/admin/menus/${id}`, {
        method: 'DELETE',
    }).then(response => {
        if (response.ok) {
            window.location.href = '/admin/menus';
        }
    });
}

function moveItem(id, direction) {
    fetch(`/admin/menus/${id}/move/${direction}`, {
        method: 'POST',
    }).then(response => {
        if (response.ok) {
            window.location.reload();
        }
    });
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
}

// static/js/admin.js - Update tag handling
function initializeTags() {
    const tagInput = document.getElementById('tag-input');
    if (!tagInput) return;

    const tagSuggestions = document.getElementById('tag-suggestions');
    const selectedTags = document.getElementById('selected-tags');
    const tagsHidden = document.getElementById('tags-hidden');
    let existingTags = [];
    let selectedTagsList = [];

    // Initialize selected tags if editing
    const initialValue = tagsHidden.value.trim();
    if (initialValue) {
        selectedTagsList = initialValue.split(',');
        updateTags();
    }

    // Fetch existing tags
    fetch('/admin/api/tags')
        .then(res => res.json())
        .then(tags => {
            existingTags = tags;
        });

    tagInput.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            const value = tagInput.value.trim();
            if (value) {
                addTag(value);
            }
        }
    });

    tagInput.addEventListener('input', () => {
        const value = tagInput.value.toLowerCase();
        if (value.length < 2) {
            tagSuggestions.style.display = 'none';
            return;
        }

        const matches = existingTags.filter(tag => 
            tag.name.toLowerCase().includes(value)
        );

        tagSuggestions.innerHTML = matches
            .map(tag => `<div class="tag-suggestion">${tag.name}</div>`)
            .join('');
        tagSuggestions.style.display = matches.length ? 'block' : 'none';
    });

    tagSuggestions.addEventListener('click', (e) => {
        if (e.target.classList.contains('tag-suggestion')) {
            addTag(e.target.textContent);
        }
    });

    function addTag(name) {
        if (!selectedTagsList.includes(name)) {
            selectedTagsList.push(name);
            updateTags();
        }
        tagInput.value = '';
        tagSuggestions.style.display = 'none';
    }

    function updateTags() {
        selectedTags.innerHTML = selectedTagsList
            .map(tag => `
                <span class="selected-tag">
                    ${tag}
                    <span class="remove-tag" data-tag="${tag}">&times;</span>
                </span>
            `).join('');
        tagsHidden.value = JSON.stringify(selectedTagsList);
    }

    selectedTags.addEventListener('click', (e) => {
        if (e.target.classList.contains('remove-tag')) {
            const tag = e.target.dataset.tag;
            selectedTagsList = selectedTagsList.filter(t => t !== tag);
            updateTags();
        }
    });
}

// Unified slug generation function
function generateSlug(text) {
    return text
        .toLowerCase()
        .trim()
        .replace(/[^a-z0-9]+/g, '-')
        .replace(/^-+|-+$/g, '');
}

// Initialize slug handling for both create and edit forms
function initializeSlugHandling() {
    const titleInput = document.getElementById('title');
    const slugInput = document.getElementById('slug');
    const form = document.querySelector('form');

    if (!titleInput || !slugInput || !form) return;

    const isEditForm = form.id === 'edit-post-form' || form.id === 'edit-page-form';
    const originalSlug = slugInput.value;

    console.log(isEditForm);

    // Auto-generate slug from title only in create forms
    if (!isEditForm) {
        titleInput.addEventListener('input', () => {
            slugInput.value = generateSlug(titleInput.value);
        });
    }

    // Show warning when slug is modified in edit forms
    if (isEditForm) {
        slugInput.addEventListener('input', () => {
            const warning = document.getElementById('slug-warning');
            const hasChanged = slugInput.value !== originalSlug;

            if (hasChanged) {
                if (!warning) {
                    const warningDiv = document.createElement('div');
                    warningDiv.id = 'slug-warning';
                    warningDiv.className = 'warning-message';
                    warningDiv.textContent = 'Warning: Changing the slug will break existing links to this content';
                    slugInput.parentNode.appendChild(warningDiv);
                }
            } else if (warning) {
                warning.remove();
            }
        });
    }

    // Validate slug format for both forms
    slugInput.addEventListener('input', () => {
        const isValidSlug = /^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(slugInput.value);
        if (!isValidSlug && slugInput.value) {
            slugInput.setCustomValidity('Slug can only contain lowercase letters, numbers, and hyphens. It cannot start or end with a hyphen.');
        } else {
            slugInput.setCustomValidity('');
        }
    });
}

function initializeTheme() {
    const html = document.documentElement;
    const theme = getCookie('admin_theme') || 'light';
    html.setAttribute('data-theme', theme);
    updateToggleIcon(theme);
}

function updateToggleIcon(theme) {
    const toggle = document.getElementById('theme-toggle');
    if (!toggle) return;

    if (theme === 'dark') {
        toggle.querySelector('.light-icon').style.display = 'none';
        toggle.querySelector('.dark-icon').style.display = 'inline';
    } else {
        toggle.querySelector('.light-icon').style.display = 'inline';
        toggle.querySelector('.dark-icon').style.display = 'none';
    }
}

function toggleTheme() {
    const html = document.documentElement;
    let currentTheme = html.getAttribute('data-theme');
    let newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    html.setAttribute('data-theme', newTheme);
    updateToggleIcon(newTheme);
    setCookie('admin_theme', newTheme, 365);
    
    // Optionally, notify the server about the theme change
    fetch('/admin/preferences', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ theme: newTheme }),
    });
}

function toggleMenu() {
    const menu = document.querySelector('.admin-nav');
    menu.classList.toggle('active');
}

// Cookie Helpers
function setCookie(name, value, days) {
    let expires = "";
    if (days) {
        const date = new Date();
        date.setTime(date.getTime() + days*24*60*60*1000);
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = name + "=" + (value || "")  + expires + "; path=/";
}

function getCookie(name) {
    const nameEQ = name + "=";
    const ca = document.cookie.split(';');
    for(let i=0;i < ca.length;i++) {
        let c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1,c.length);
        if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
    }
    return null;
}

function initializePublishDateToggle() {
    const publishType = document.getElementById('publishType');
    const publishDateGroup = document.getElementById('publishDateGroup');
    const publishedAtInput = document.getElementById('publishedAt');
    
    if (!publishType || !publishDateGroup || !publishedAtInput) return;

    function updatePublishDate() {
        if (publishType.value === 'immediately') {
            publishDateGroup.style.display = 'none';
            publishedAtInput.removeAttribute('required');
        } else {
            publishDateGroup.style.display = 'block';
            publishedAtInput.setAttribute('required', 'required');
        }
    }

    publishType.addEventListener('change', updatePublishDate);
    
    // Initial state
    if (publishedAtInput.value) {
        publishType.value = 'scheduled';
    }
    updatePublishDate();
}

function initializeMenuItemForm() {
    const pageSelect = document.getElementById('page_id');
    const urlInput = document.getElementById('url');
    const labelInput = document.getElementById('label');
    const form = document.querySelector('form');

    if (!pageSelect || !urlInput || !labelInput || !form) return;

    pageSelect.addEventListener('change', function() {
        const pageId = this.value;
        const selectedOption = this.options[this.selectedIndex];
        
        if (pageId) {
            const pageSlug = selectedOption.getAttribute('data-slug');
            urlInput.value = '/pages/' + pageSlug;
            if (!labelInput.value) {
                labelInput.value = selectedOption.text;
            }
            urlInput.readOnly = true;
        } else {
            urlInput.readOnly = false;
        }
    });

    // Clear page selection when URL is manually edited
    urlInput.addEventListener('input', function() {
        if (this.value !== this.defaultValue) {
            pageSelect.value = '';
            this.readOnly = false;
        }
    });

    form.addEventListener('submit', function(e) {
        e.preventDefault();
        
        // Validate form
        if (!labelInput.value.trim()) {
            alert('Please enter a label for the menu item');
            return;
        }

        if (!urlInput.value.trim() && !pageSelect.value) {
            alert('Please either enter a URL or select a page');
            return;
        }

        // If validation passes, submit the form
        this.submit();
    });
}

function initializeMenuItems() {
    const moveUpButtons = document.querySelectorAll('.move-up');
    const moveDownButtons = document.querySelectorAll('.move-down');
    const deleteMenuButtons = document.querySelectorAll('.delete-menu-item');

    moveUpButtons.forEach(button => {
        button.addEventListener('click', function() {
            const id = this.getAttribute('data-id');
            moveItem(id, 'up');
        });
    });

    moveDownButtons.forEach(button => {
        button.addEventListener('click', function() {
            const id = this.getAttribute('data-id');
            moveItem(id, 'down');
        });
    });

    deleteMenuButtons.forEach(button => {
        button.addEventListener('click', function() {
            const id = this.getAttribute('data-id');
            deleteMenuItem(id);
        });
    });
}

// Initialize on DOM Content Loaded
document.addEventListener('DOMContentLoaded', () => {
    initializeEditor();
    initializeTags();
    initializeSlugHandling();
    initializeTheme();
    initializeMenuItemForm();
    initializeMenuItems();
    initializePublishDateToggle();

    const themeToggle = document.getElementById('theme-toggle');
    if (themeToggle) {
        themeToggle.addEventListener('click', toggleTheme);
    }

    const menuToggle = document.getElementById('menu-toggle');
    if (menuToggle) {
        menuToggle.addEventListener('click', toggleMenu);
    }
});