
.PHONY: doc
doc:
	-$(RM) doc/*.1
	mkdir -p doc
	cd doc && go run "../hack/generate-man.go"
