.PHONY: install_microk8s

install_microk8s:
	uvx --from=ansible ansible-playbook -i inventory/hosts install_microk8s.yml