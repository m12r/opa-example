package demo.api

import future.keywords.if

test_allow_guest_user_to_read_homepage if {
    allow with input as {
        "method": "GET",
        "path": [""],
        "user": ""
    }
}

test_allow_foo_user_to_read_homepage if {
    allow with input as {
        "method": "GET",
        "path": [""],
        "user": "foo"
    }
}

test_allow_guest_user_to_read_health if {
    allow with input as {
        "method": "GET",
        "path": ["health"],
        "user": ""
    }
}

test_allow_foo_user_to_read_health if {
    allow with input as {
        "method": "GET",
        "path": ["health"],
        "user": "foo"
    }
}

test_deny_guest_user_to_read_foo_user_payments if {
    not allow with input as {
        "method": "GET",
        "path": ["payments", "foo"],
        "user": ""
    }
}

test_deny_bar_user_to_read_foo_payments if {
    not allow with input as {
        "method": "GET",
        "path": ["payments", "foo"],
        "user": "bar"
    }
}

test_allow_foo_user_to_read_foo_payments if {
    allow with input as {
        "method": "GET",
        "path": ["payments", "foo"],
        "user": "foo"
    }
}

test_allow_teller_user_to_read_foo_payments if {
    allow with input as {
        "method": "GET",
        "path": ["payments", "foo"],
        "user": "teller"
    }
}
