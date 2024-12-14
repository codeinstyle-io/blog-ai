function deleteTag(id) {
    fetch(`/admin/tags/${id}`, {
        method: 'DELETE',
    }).then((response) => {
        if (response.ok) {
            window.location.href = '/admin/tags';
        } else {
            response.json().then(data => {
                showError(data.error || 'Failed to delete tag');
            }).catch(() => {
                showError('Failed to delete tag');
            });
        }
    }).catch(error => {
        showError('Failed to delete tag');
        console.error('Error:', error);
    });
}

function deletePost(id) {
    fetch(`/admin/posts/${id}`, {
        method: 'DELETE',
    }).then((response) => {
        if (response.ok) {
            window.location.href = '/admin/posts';
        } else {
            response.json().then(data => {
                showError(data.error || 'Failed to delete post');
            }).catch(() => {
                showError('Failed to delete post');
            });
        }
    }).catch(error => {
        showError('Failed to delete post');
        console.error('Error:', error);
    });
}

function deletePage(id) {
    fetch(`/admin/pages/${id}`, {
        method: 'DELETE',
    }).then(response => {
        if (response.ok) {
            window.location.href = '/admin/pages';
        } else {
            response.json().then(data => {
                showError(data.error || 'Failed to delete page');
            }).catch(() => {
                showError('Failed to delete page');
            });
        }
    }).catch(error => {
        showError('Failed to delete page');
        console.error('Error:', error);
    });
}

function deleteMenuItem(id) {
    fetch(`/admin/menus/${id}`, {
        method: 'DELETE',
    }).then(response => {
        if (response.ok) {
            window.location.href = '/admin/menus';
        } else {
            response.json().then(data => {
                showError(data.error || 'Failed to delete menu item');
            }).catch(() => {
                showError('Failed to delete menu item');
            });
        }
    }).catch(error => {
        showError('Failed to delete menu item');
        console.error('Error:', error);
    });
}

function deleteMedia(id) {
    fetch(`/admin/media/${id}`, {
        method: 'DELETE',
    }).then((response) => {
        if (response.ok) {
            window.location.href = '/admin/media';
        } else {
            response.json().then(data => {
                showError(data.error || 'Failed to delete media');
            }).catch(() => {
                showError('Failed to delete media');
            });
        }
    }).catch(error => {
        showError('Failed to delete media');
        console.error('Error:', error);
    });
}

function showError(message) {
    const errorDiv = document.createElement('div');
    errorDiv.className = 'alert alert-error';
    errorDiv.textContent = message;
    
    // Insert at the top of the main content area
    const mainContent = document.querySelector('main');
    if (mainContent) {
        mainContent.insertBefore(errorDiv, mainContent.firstChild);
    }
    
    // Remove after 5 seconds
    setTimeout(() => {
        errorDiv.remove();
    }, 5000);
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
        tagsHidden.value = selectedTagsList.join(',');  
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

    // Initialize form fields from server data
    const initialPageId = pageSelect.value;
    if (initialPageId) {
        const selectedOption = pageSelect.options[pageSelect.selectedIndex];
        if (selectedOption) {
            const pageSlug = selectedOption.getAttribute('data-slug');
            urlInput.value = '/pages/' + pageSlug;
            urlInput.readOnly = true;
        }
    }

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
        // Validate form
        if (!labelInput.value.trim()) {
            e.preventDefault();
            alert('Please enter a label for the menu item');
            return;
        }

        if (!urlInput.value.trim() && !pageSelect.value) {
            e.preventDefault();
            alert('Please either enter a URL or select a page');
            return;
        }

        // Form is valid - let it submit naturally
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

function moveItem(id, direction) {
    const button = document.querySelector(`button[data-id="${id}"].move-${direction}`);
    if (button && button.disabled) {
        return; // Don't move if button is disabled
    }

    fetch(`/admin/menus/${id}/move/${direction}`, {
        method: 'POST',
    }).then(response => {
        if (response.ok) {
            window.location.reload();
        } else {
            response.json().then(data => {
                showError(data.error || 'Failed to move item');
            }).catch(() => {
                showError('Failed to move item');
            });
        }
    }).catch(error => {
        showError('Failed to move item');
        console.error('Error:', error);
    });
}

function initializeMenuToggle() {
    const menuToggle = document.getElementById('menu-toggle');
    const adminNav = document.querySelector('.admin-nav');
    
    if (menuToggle && adminNav) {
        menuToggle.innerHTML = '<i class="fas fa-bars"></i>';
        
        menuToggle.addEventListener('click', () => {
            adminNav.classList.toggle('active');
            // Update aria-expanded for accessibility
            const isExpanded = adminNav.classList.contains('active');
            menuToggle.setAttribute('aria-expanded', isExpanded);
        });

        // Close menu when clicking outside
        document.addEventListener('click', (event) => {
            if (!adminNav.contains(event.target) && 
                !menuToggle.contains(event.target) && 
                adminNav.classList.contains('active')) {
                adminNav.classList.remove('active');
                menuToggle.setAttribute('aria-expanded', 'false');
            }
        });
    }
}

// Initialize on DOM Content Loaded
document.addEventListener('DOMContentLoaded', () => {
    initializeEditor();
    initializeTags();
    initializeSlugHandling();
    initializePublishDateToggle();
    initializeMenuItemForm();
    initializeMenuItems();
    initializeMenuToggle();
});
