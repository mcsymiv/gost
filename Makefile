.PHONY: home driver rec harvest tc

# parses and passes arguments to "make home @args"
# Ex.: 
# 	make home hello
# In code:
# 	fmt.Println(os.Args)
# 	0: test binary path
# 	1: -test.paniconexit0
# 	2: -test.timeout=10m0s
# 	3: -test.v=true
# 	4: -test.count=1
# 	5: -test.run=TestHome
# 	6: hello
# If the first argument is "run"...
ifneq (,$(filter $(firstword $(MAKECMDGOALS)), home harvest tc))
# ifeq (home,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  args := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(args):;@:)
endif

home:
	go test -v -count=1 test/home_test.go -run TestHome $(args)

rec:
	go test -v -count=1 test/record_test.go -run TestRecord

harvest:
	go test -v -count=1 test/harvest_test.go -run TestHarvest $(args)

tc:
	go test -v -count=1 test/tc_test.go -run TestTc $(args)
