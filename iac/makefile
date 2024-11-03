.PHONY: install_microk8s

encrypt_secrets:
	uvx --from=ansible ansible-vault encrypt secrets.yml

install_microk8s:
	uvx --from=ansible ansible-playbook -i inventory/hosts -e @secrets.yml --ask-vault-pass install_microk8s.yml

uninstall_microk8s:
	uvx --from=ansible ansible-playbook -i inventory/hosts uninstall_microk8s.yml

install_blog:
	uvx --from=ansible ansible-playbook -i inventory/hosts install_blog.yml