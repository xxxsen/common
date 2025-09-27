package ipcheck

import "testing"

func TestCheckerAddAndCheckIPv4(t *testing.T) {
	var checker Checker

	if err := checker.Add("192.168.1.10"); err != nil {
		t.Fatalf("Add IPv4 failed: %v", err)
	}

	match, err := checker.Check("192.168.1.10")
	if err != nil || !match {
		t.Fatalf("expected match for exact IPv4, got match=%v err=%v", match, err)
	}

	match, err = checker.Check("192.168.1.11")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if match {
		t.Fatal("unexpected match for different IPv4")
	}
}

func TestCheckerAddAndCheckIPv6(t *testing.T) {
	var checker Checker

	ipv6 := "2001:db8::1"
	if err := checker.Add(ipv6); err != nil {
		t.Fatalf("Add IPv6 failed: %v", err)
	}

	match, err := checker.Check(ipv6)
	if err != nil || !match {
		t.Fatalf("expected match for exact IPv6, got match=%v err=%v", match, err)
	}

	match, err = checker.Check("2001:db8::2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if match {
		t.Fatal("unexpected match for different IPv6")
	}
}

func TestCheckerCIDR(t *testing.T) {
	var checker Checker

	if err := checker.Add("10.0.0.0/24"); err != nil {
		t.Fatalf("Add CIDR failed: %v", err)
	}

	match, err := checker.Check("10.0.0.5")
	if err != nil || !match {
		t.Fatalf("expected match in CIDR, got match=%v err=%v", match, err)
	}

	match, err = checker.Check("10.0.1.5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if match {
		t.Fatal("unexpected match outside CIDR")
	}
}

func TestCheckerCIDRIPv6(t *testing.T) {
	var checker Checker

	if err := checker.Add("2001:db8::/126"); err != nil {
		t.Fatalf("Add IPv6 CIDR failed: %v", err)
	}

	match, err := checker.Check("2001:db8::1")
	if err != nil || !match {
		t.Fatalf("expected IPv6 match, got match=%v err=%v", match, err)
	}

	match, err = checker.Check("2001:db8::5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if match {
		t.Fatal("unexpected match outside IPv6 CIDR")
	}
}

func TestCheckerInvalidInput(t *testing.T) {
	var checker Checker

	if err := checker.Add("not-an-ip"); err == nil {
		t.Fatal("expected error for invalid IP")
	}

	if err := checker.Add("10.0.0.0/33"); err == nil {
		t.Fatal("expected error for invalid CIDR")
	}

	if _, err := checker.Check("bad-ip"); err == nil {
		t.Fatal("expected error for invalid check IP")
	}
}

func TestCheckerTrimmedInput(t *testing.T) {
	var checker Checker
	if err := checker.Add("  172.16.0.1  "); err != nil {
		t.Fatalf("Add trimmed IP failed: %v", err)
	}

	match, err := checker.Check("\t172.16.0.1\n")
	if err != nil || !match {
		t.Fatalf("expected match for trimmed input, got match=%v err=%v", match, err)
	}
}
