#!/usr/bin/env ruby

require 'curb'
require 'json'

products = [
  # productID, registerID

#   {pid: "0797984e-548b-4d4d-b33e-98461e72fe90", rid: "2b1a2fe3-b7ef-41a6-86fa-d4eb076c8a24"},
#   {pid: "7c3e83d2-fbb3-454d-92ae-a2d756886e41", rid: "02c93e2f-d390-4dd1-8baf-5458ae3af9e3"},

# p20-p40
  {pid: "b9f3fab4-bae4-4e23-9af4-d08046f12f28", rid: "e06e629b-c93c-4134-b75d-345cf2cde098"},
  {pid: "ac7b69a5-96e0-408e-9f09-3f751c76fd6a", rid: "8bf16a31-5483-4e2d-88dc-c12208b4ef56"},
  {pid: "9938393f-9b70-49b5-abcf-3e1c430492b4", rid: "1f749e3c-ab85-47c8-a5c4-e8a5abb2a36a"},
  {pid: "c7839062-af85-426a-a6eb-5756fc2b4619", rid: "e1b8c31d-8479-435b-a184-dd826792e1da"},
  {pid: "4b351d4c-7bfe-4521-ada5-f79b230ecfe8", rid: "aae9d9a0-b585-46c6-8323-1785168e2c50"},
  {pid: "b5e4d7ab-e2d9-4bc2-8c71-6e9b345bb9c4", rid: "43055388-eed3-43d4-85d8-14ecd60ac992"},
  {pid: "e92e84b4-27c8-4372-93b6-09531089e116", rid: "da95ee7b-a292-4b0b-8966-26bfdd9f93e9"},
  {pid: "272da82b-e40e-401d-9216-aa7d33f0cce8", rid: "3302ae4b-d942-417e-8230-cffb633e7d5c"},
  {pid: "1942337f-4354-4a89-a34c-1d34e6567c49", rid: "e6982406-9779-4724-a5e2-48f306630527"},
  {pid: "29a95cb7-d638-4f26-b454-0848a17cbb2a", rid: "9d6937d1-909c-4122-af4c-1e92dee90d7c"},
  {pid: "f6ee2f56-18e0-4867-8de4-a25e8d4cd6e2", rid: "447d197c-5eb1-4c9b-9b42-bcb8209ecc34"},
  {pid: "cd6072bb-dbb8-43cb-b6c2-bbf94654906d", rid: "dc8f2e8b-e0dd-4742-ace8-893301b066a4"},
  {pid: "d515ca7f-9f68-46bd-a3b1-d884adad8d28", rid: "f25e40d7-4a1a-4e98-9c79-57a0b4b62729"},
  {pid: "47db929f-c2de-43d3-a59d-4940bf2c8677", rid: "b5f83955-4dc5-430d-8e08-0c13f9ecff50"},
  {pid: "85db38b9-0f59-4975-85fd-8328919cd8e8", rid: "ef88407e-2d3e-489d-9469-0ecb45d5a1bc"},
  {pid: "409305a8-3218-4787-aec3-e6a36a1c7c69", rid: "812009db-9ee7-4d33-aa7f-95bfa957b12a"},
  {pid: "1f95bd13-4330-4ba2-b6e2-89200cb3b64b", rid: "1584a766-dbdd-4ac6-b009-c96de677fe94"},
  {pid: "f8b892ae-0b37-44c1-a355-4bc0f650141d", rid: "4ce0ca67-cda7-48e8-993d-c4bce47ed643"},
  {pid: "cdf3af78-9a89-4811-8671-0c1ebdb68499", rid: "29ab5101-7aea-4357-b3a9-a43c9149357f"},
  {pid: "811c9ac8-e401-44fa-9b09-a5386ad857ed", rid: "b010ce0e-7c84-46ad-8e6a-b3a0a575c66b"},
  {pid: "f9a47379-5251-453c-b1ce-34b725e67420", rid: "e5a7f8b7-b8e5-4f46-a440-ecb35a83fa68"},
]

# puts products[0][:pid]


def getClose(productID, steakID)
  host = "localhost:8080"
  c = Curl::Easy.perform(host+'/api/steak/detail?productID='+productID+'&steakID='+steakID)

  j = JSON.parse(c.body_str)
#   puts j['closeID']

  c = {
    steakID: steakID,
    closeID: j['closeID'],
    weChatID: "big-eater"
  }
  puts c.to_json

  http = Curl.post(host+'/api/steak/close', c.to_json)
  puts http.body_str
end

products.each { |product|
  pid = product[:pid]
  rid = product[:rid]

  getClose(pid, rid)
}