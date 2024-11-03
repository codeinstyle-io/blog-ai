.PHONY: install_microk8s

encrypt_secrets:
	uvx --from=ansible ansible-vault encrypt secrets.yml

install_microk8s:
	uvx --from=ansible ansible-playbook -e @secrets.yml --ask-vault-pass -i inventory/hosts install_microk8s.yml

uninstall_microk8s:
	uvx --from=ansible ansible-playbook -i inventory/hosts uninstall_microk8s.yml