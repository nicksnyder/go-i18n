all:
	$(MAKE) -C goi18nextract/
	$(MAKE) -C i18n/

clean:
	$(MAKE) -C goi18nextract/ clean
	$(MAKE) -C i18n/ clean

install:
	$(MAKE) -C goi18nextract/ install
	$(MAKE) -C i18n/ install
