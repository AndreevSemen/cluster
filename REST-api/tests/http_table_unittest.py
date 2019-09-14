import requests
import json
import unittest
import copy

server_addr = 'http://localhost:80'
get_url = server_addr + '/gossip/table'
drop_url = get_url + '/drop'
update_url = get_url + '/update'
servers_infos = [{"ip": "127.0.0.1", "port": 80, "state": "alive"},
                 {"ip": "127.0.0.1", "port": 81, "state": "dead"},
                 {"ip": "127.0.0.1", "port": 82, "state": "onload"},
                 {"ip": "127.0.0.2", "port": 80, "state": "dead"},
                 {"ip": "127.0.0.3", "port": 81, "state": "inactive"}]


def get():
    return json.loads(requests.get(get_url).content)


def post(input_json):
    return json.loads(requests.post(update_url, json=input_json).content)


def drop_database():
    requests.delete(drop_url)


def empty(arr):
    return arr is None


def ips_equal_to(arr, ip):
    for server in arr:
        if server["ip"] != ip:
            return False

    return True


class RestAPITableCase(unittest.TestCase):
    def test_empty(self):
        drop_database()

        self.assertTrue(empty(get()["servers"]))

    def test_post_to_empty(self):
        drop_database()

        curr_json = servers_infos[0]
        self.assertTrue(post(curr_json)["servers"])

    def test_post_neighbours(self):
        drop_database()

        curr_json = servers_infos[0] # 127.0.0.1:80
        post(curr_json)
        curr_json = servers_infos[1] # 127.0.0.1:81
        resp_json = post(curr_json)
        self.assertEqual(len(resp_json["servers"]), 1)
        self.assertTrue(ips_equal_to(resp_json["servers"], curr_json["ip"]))

        curr_json = servers_infos[2] # 127.0.0.1:82
        resp_json = post(curr_json)
        self.assertEqual(len(resp_json["servers"]), 2)
        self.assertTrue(ips_equal_to(resp_json["servers"], curr_json["ip"]))

        curr_json = servers_infos[3] # 127.0.0.2:80
        resp_json = post(curr_json)
        self.assertTrue(empty(resp_json["servers"]))

    def test_update(self):
        drop_database()

        old_json = copy.deepcopy(servers_infos[0])
        post(old_json)

        new_json = copy.deepcopy(servers_infos[0])
        new_json["state"] = "inactive"

        resp_json = get()
        self.assertTrue(old_json in resp_json["servers"])
        self.assertTrue(new_json not in resp_json["servers"])

        post(new_json)
        resp_json = get()
        self.assertTrue(old_json not in resp_json["servers"])
        self.assertTrue(new_json in resp_json["servers"])

    def test_incorrect_state(self):
        drop_database()

        incorrect_json = copy.deepcopy(servers_infos[0])
        incorrect_json["state"] = "semi-alive"

        resp = requests.post(update_url, json=incorrect_json)
        self.assertEqual(resp.content, b'sent server with invalid state')


if __name__ == '__main__':
    unittest.main()
