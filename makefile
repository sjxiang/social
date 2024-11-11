

run:
	@echo "运行"
	go run cmd/api/*.go


mysql-container-console:
	@echo "MySQL 容器控制台"
	@echo "mysql --host=127.0.0.1 --port=3306 --user=root --password=my-secret-pw"
	docker exec -it db sh


redis-container-console:
	@echo "Redis 容器控制台"
	@echo "redis-cli"
	docker exec -it cache sh


.PHONY: run mysql-container-console redis-container-console

