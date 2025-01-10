function deleteTag(id) {
    fetch(`/admin/tags/${id}`, {
        method: 'DELETE',
    }).then((response) => response.json())
    .then((data) => {
        if (data.redirect) {
            window.location.href = data.redirect;
        }
    }).catch(error => {
        console.error('Error:', error);
    });
}

function deletePost(id) {
    fetch(`/admin/posts/${id}`, {
        method: 'DELETE',
    }).then((response) => response.json())
    .then((data) => {
        if (data.redirect) {
            window.location.href = data.redirect;
        }
    }).catch(error => {
        console.error('Error:', error);
    });
}

function deletePage(id) {
    fetch(`/admin/pages/${id}`, {
        method: 'DELETE',
    }).then((response) => response.json())
    .then((data) => {
        if (data.redirect) {
            window.location.href = data.redirect;
        }
    }).catch(error => {
        console.error('Error:', error);
    });
}

function deleteMenuItem(id) {
    fetch(`/admin/menus/${id}`, {
        method: 'DELETE',
    }).then((response) => response.json())
    .then((data) => {
        if (data.redirect) {
            window.location.href = data.redirect;
        }
    }).catch(error => {
        console.error('Error:', error);
    });
}

function deleteMedia(id) {
    fetch(`/admin/media/${id}`, {
        method: 'DELETE',
    }).then((response) => response.json())
    .then((data) => {
        if (data.redirect) {
            window.location.href = data.redirect;
        }
    }).catch(error => {
        console.error('Error:', error);
    });
}

function deleteUser(id) {
    fetch(`/admin/users/${id}`, {
        method: 'DELETE',
    }).then((response) => response.json())
    .then((data) => {
        if (data.redirect) {
            window.location.href = data.redirect;
        }
    }).catch(error => {
        console.error('Error:', error);
    });
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
                console.error(data.error || 'Failed to move item');
            }).catch(() => {
                console.error('Failed to move item');
            });
        }
    }).catch(error => {
        console.error('Failed to move item');
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

!(function(win, doc) {
    function openMediaModal(cb) {
        doc.getElementById('mediaModal').style.display = 'block';
        loadMediaItems(cb);
    }

    function closeMediaModal() {
        doc.getElementById('mediaModal').style.display = 'none';
    }

    function loadMediaItems(cb) {
        fetch('/admin/api/media')
            .then(response => response.json())
            .then(items => {
                mediaItems = items;
                const grid =    doc.getElementById('mediaGrid');
                grid.innerHTML = '';

                items.forEach(item => {
                    const div = doc.createElement('div');
                    div.className = 'media-item';

                    if (item.MimeType.startsWith('image/')) {
                        div.innerHTML = `
                            <div class="media-preview">
                                <img src="/media/${item.Path}" alt="${item.Name}">
                            </div>
                            <div class="media-info">
                                <h3>${item.Name}</h3>
                            </div>
                        `;
                    } else {
                        div.innerHTML = `
                            <div class="media-preview file">
                                <i class="fas fa-file"></i>
                            </div>
                            <div class="media-info">
                                <h3>${item.Name}</h3>
                            </div>
                        `;
                    }

                    div.onclick = () => cb(item);
                    grid.appendChild(div);
                });
            })
            .catch(error => console.error('Error loading media:', error));
    }

    function insertMedia(media, currentEditorId) {
        const editor = doc.getElementById(currentEditorId);
        const format = doc.querySelector('input[name="format"]:checked').value;
        let tag;

        if (media.MimeType.startsWith('image/')) {
            // For images
            tag = format === 'markdown'
                ? `![${media.Name}](/media/${media.Path})`
                : `<img src="/media/${media.Path}" alt="${media.Name}">`;
        } else {
            // For other files
            tag = format === 'markdown'
                ? `[${media.Name}](/media/${media.Path})`
                : `<a href="/media/${media.Path}">${media.Name}</a>`;
        }

        // Get cursor position
        const start = editor.selectionStart;
        const end = editor.selectionEnd;

        // Insert the tag at cursor position
        editor.value = editor.value.substring(0, start) + tag + editor.value.substring(end);
    }

    win.openMediaModal = openMediaModal;
    win.closeMediaModal = closeMediaModal;
    win.insertMedia = insertMedia;

    // Close modal when clicking outside
    win.onclick = function(event) {
        const modal = doc.getElementById('mediaModal');
        if (event.target == modal) {
            closeMediaModal();
        }
    };
}(window, document));

function openEditorMediaSelector(editorId) {
    openMediaModal((media) => {
        insertMedia(media, editorId);
        closeMediaModal();
    });
}

function openLogoMediaSelector() {
    openMediaModal((media) => {
        const logoInput = document.getElementById('logo_id');
        const imagePreview = document.querySelector('.image-preview');
        const img = document.createElement('img');
        img.src = `/media/${media.Path}`;
        if (imagePreview) {
            imagePreview.replaceChildren(img);
        }
        logoInput.value = media.ID;
        closeMediaModal();
    });
}

(function() {

    document.querySelectorAll('[x-dynamic-date]').forEach((element) => {
        const datetime = element.getAttribute('x-dynamic-date');
        const timezone = element.getAttribute('x-dynamic-date-timezone');

        // TODO: Localize
        const date = new Date(datetime).toLocaleString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });

        element.innerHTML = date;
    });

    Inity.register('posts', Apps.Posts, { onSubmit: async(data, done, error, props) => {
        let method = 'POST';
        let url = '/admin/api/posts';

        if(props.id) {
            method = 'PUT';
            url = url + '/' + props.id;
        }
        done('saving');

        const resp = await fetch(url, {
        method,
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
        });
        const json = await resp.json();

        if(resp.ok) {
            if (json.redirect) {
                done('saved');
                window.location.href = json.redirect;
            }
        } else {
            error(json.error);
            document.querySelector('.editor-container').scrollIntoView()
        }
    }
    });

    Inity.register('pages', Apps.Pages, { onSubmit: async(data, done, error, props) => {
        let method = 'POST';
        let url = '/admin/api/pages';

        if(props.id) {
        method = 'PUT';
        url = url + '/' + props.id;
        }
        done('saving');

        const resp = await fetch(url, {
        method,
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
        });

        if(resp.ok) {
            const json = await resp.json();
            if (json.redirect) {
                done('saved');
                window.location.href = json.redirect;
            }
        } else {
            error(json.error);
            document.querySelector('.editor-container').scrollIntoView()
        }
    }
    });

  document.addEventListener("DOMContentLoaded", () => Inity.attach());
})();

// Initialize on DOM Content Loaded
document.addEventListener('DOMContentLoaded', () => {
    initializeMenuItemForm();
    initializeMenuItems();
    initializeMenuToggle();
});
