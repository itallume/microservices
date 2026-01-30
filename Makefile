.PHONY: build-order build-payment build-shipping build-all
.PHONY: deploy-order deploy-payment deploy-shipping deploy-all deploy-mysql
.PHONY: delete-order delete-payment delete-shipping delete-all delete-mysql
.PHONY: logs-order logs-payment logs-shipping logs-mysql
.PHONY: status restart-all

# ==================== BUILD ====================
build-order:
	cd .. && docker build -t order -f microservices/order/Dockerfile .

build-payment:
	cd .. && docker build -t payment -f microservices/payment/Dockerfile .

build-shipping:
	cd .. && docker build -t shipping -f microservices/shipping/Dockerfile .

build-all: build-order build-payment build-shipping
	@echo "All images built successfully!"

# ==================== DEPLOY ====================
deploy-mysql:
	kubectl apply -f mysql/deployment.yaml

deploy-order:
	kubectl apply -f order/deployment.yaml

deploy-payment:
	kubectl apply -f payment/deployment.yaml

deploy-shipping:
	kubectl apply -f shipping/deployment.yaml

deploy-all: deploy-mysql deploy-order deploy-payment deploy-shipping
	@echo "All services deployed!"

# ==================== DELETE ====================
delete-mysql:
	kubectl delete -f mysql/deployment.yaml

delete-order:
	kubectl delete -f order/deployment.yaml

delete-payment:
	kubectl delete -f payment/deployment.yaml

delete-shipping:
	kubectl delete -f shipping/deployment.yaml

delete-all: delete-order delete-payment delete-shipping delete-mysql
	@echo "All services deleted!"

# ==================== LOGS ====================
logs-mysql:
	kubectl logs -f -l service=mysql

logs-order:
	kubectl logs -f -l service=order

logs-payment:
	kubectl logs -f -l service=payment

logs-shipping:
	kubectl logs -f -l service=shipping

# ==================== STATUS ====================
status:
	@echo "=== PODS ==="
	kubectl get pods
	@echo ""
	@echo "=== SERVICES ==="
	kubectl get services
	@echo ""
	@echo "=== INGRESS ==="
	kubectl get ingress

# ==================== RESTART ====================
restart-order:
	kubectl rollout restart deployment/order

restart-payment:
	kubectl rollout restart deployment/payment

restart-shipping:
	kubectl rollout restart deployment/shipping

restart-all: restart-order restart-payment restart-shipping
	@echo "All services restarted!"

# ==================== FULL DEPLOY ====================
up: build-all deploy-all status
	@echo "All microservices are up!"

down: delete-all
	@echo "All microservices are down!"
