import requests
import json
import unittest
import copy

addr = 'http://localhost:84'
get_url = addr + '/gossip/table'
update_url = get_url + '/update'
servers_infos = [{"ip": "127.0.0.1", "port": 80, "state": "alive"},
                 {"ip": "127.0.0.1", "port": 81, "state": "dead"},
                 {"ip": "127.0.0.1", "port": 82, "state": "onload"},
                 {"ip": "127.0.0.2", "port": 80, "state": "dead"},
                 {"ip": "127.0.0.3", "port": 81, "state": "inactive"}]


def get_result():
    return json.loads(requests.get(get_url).content)


def post_result(in_json):
    return json.loads(requests.post(update_url, json=in_json).content)


def empty(arr):
    return arr is None


def ips_equal_to(arr, ip):
    for server in arr:
        if server["ip"] != ip:
            return False

    return True


class RestAPITest(unittest.TestCase):
    def empty_table(self):
        self.assertTrue(empty(get_result()["servers"]))

    def post_to_empty(self):
        curr_json = servers_infos[0]
        self.assertTrue(empty(post_result(curr_json)["servers"]))

    def post_neighbours(self):
        curr_json = servers_infos[1]
        resp_json = post_result(curr_json)
        self.assertEqual(len(resp_json["servers"]), 1)
        self.assertTrue(ips_equal_to(resp_json["servers"], curr_json["ip"]))

        curr_json = servers_infos[2]
        resp_json = post_result(curr_json)
        self.assertEqual(len(resp_json["servers"]), 2)
        self.assertTrue(ips_equal_to(resp_json["servers"], curr_json["ip"]))

        curr_json = servers_infos[3]
        resp_json = post_result(curr_json)
        self.assertTrue(empty(resp_json["servers"]))

    def update(self):
        old_json = copy.deepcopy(servers_infos[0])
        new_json = copy.deepcopy(servers_infos[0])
        new_json["state"] = "inactive"

        resp_json = get_result()
        self.assertTrue(old_json in resp_json["servers"])
        self.assertTrue(new_json not in resp_json["servers"])

        post_result(new_json)
        resp_json = get_result()
        self.assertTrue(old_json not in resp_json["servers"])
        self.assertTrue(new_json in resp_json["servers"])

    def incorrect_state(self):
        incorrect_json = copy.deepcopy(servers_infos[0])
        incorrect_json["state"] = "semi-alive"

        with self.assertRaises(ValueError):
            resp_json = post_result(incorrect_json)


if __name__ == '__main__':
    unittest.main()
