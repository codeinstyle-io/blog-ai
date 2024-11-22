function deleteTag(id) {
    if (confirm('Are you sure you want to delete this tag?')) {
        fetch(`/admin/tags/${id}`, {
            method: 'DELETE',
        }).then(() => {
            window.location.reload();
        });
    }
}

function deletePost(id) {
    fetch(`/admin/posts/${id}`, {
        method: 'DELETE',
    }).then(() => {
        window.location.href = '/admin/posts';
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

// static/js/admin.js - Update tag handling
function initializeTags() {
    const tagInput = document.getElementById('tag-input');
    const tagSuggestions = document.getElementById('tag-suggestions');
    const selectedTags = document.getElementById('selected-tags');
    const tagsHidden = document.getElementById('tags-hidden');
    let existingTags = [];
    let selectedTagsList = [];

    // Initialize selected tags if editing
    selectedTags.querySelectorAll('.selected-tag').forEach(tag => {
        selectedTagsList.push(tag.textContent.trim());
    });
    
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

function initializeSlugWarning() {
    const slugInput = document.getElementById('slug');
    const originalSlug = slugInput.value;

    slugInput.addEventListener('input', () => {
        if (originalSlug && slugInput.value !== originalSlug) {
            if (!document.getElementById('slug-warning')) {
                const warning = document.createElement('div');
                warning.id = 'slug-warning';
                warning.className = 'warning-message';
                warning.textContent = 'Warning: Changing the slug will break existing links to this post';
                slugInput.parentNode.appendChild(warning);
            }
        } else {
            const warning = document.getElementById('slug-warning');
            if (warning) {
                warning.remove();
            }
        }
    });
}

document.addEventListener("DOMContentLoaded", () => {
    initializeEditor();
    initializeTags();
    initializeSlugWarning();
});
