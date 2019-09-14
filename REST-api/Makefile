run-tests:
	cd .. && sudo make docker && sudo make daemon
	ping -c 2 127.0.0.1
	python3 tests/http_table_unittest.py
	cd .. && sudo make stop