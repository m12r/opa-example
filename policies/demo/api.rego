package demo.api

users := {
    "teller": ["teller"]
}

default current_user = ""

current_user := user {
    user := input.user
}

user_levels[l] {
    current_user != ""
    l := users[current_user][_]
}

user_levels[l] {
    current_user != ""
    l := "authorized_user"
}

user_levels[l] {
    current_user == ""
    l := "guest_user"
}

default allow = false

operations["read_homepage"] {
    input.method == "GET"
    input.path == [""]
}

operations["read_health"] {
    input.method == "GET"
    input.path == ["health"]
}

operations["read_own_payments"] {
    input.method == "GET"
    input.path == ["payments", current_user]
}

operations["read_all_payments"] {
    input.method == "GET"
    array.slice(input.path, 0, 1) == ["payments"]
}

default can_read_public = false

default can_read_private = false

can_read_public {
    some user_level
    user_levels[user_level]
    ["guest_user", "authorized_user", "teller"][_] == user_level
}

can_read_private {
    some user_level
    user_levels[user_level]
    ["authorized_user", "teller"][_] == user_level
}

can_read_payments {
    some user_level
    user_levels[user_level]
    ["teller"][_] == user_level
}

allow {
    operations.read_homepage
    can_read_public
}

allow {
    operations.read_health
    can_read_public
}

allow {
    operations.read_own_payments
    can_read_private
}

allow {
    operations.read_all_payments
    can_read_payments
}
