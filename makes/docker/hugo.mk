HUGO_SITE ?= html
HUGO_PORT ?= 8080
HUGO_HOST ?= localhost:$(HUGO_PORT)

.PHONY: hugo-init
hugo-init:
	mkdir -p site/$(HUGO_SITE)/archetypes \
	         site/$(HUGO_SITE)/content \
	         site/$(HUGO_SITE)/data \
	         site/$(HUGO_SITE)/i18n \
	         site/$(HUGO_SITE)/layouts \
	         site/$(HUGO_SITE)/static \
	         site/$(HUGO_SITE)/themes
	for file in $$(ls site/$(HUGO_SITE)); do \
	    if [[ -d site/$(HUGO_SITE)/$$file ]]; then \
	        export COUNT=$$(ls -a site/$(HUGO_SITE)/$$file | wc -l); \
	        if [[ $$COUNT -lt 3 ]]; then \
	            touch site/$(HUGO_SITE)/$$file/.gitkeep; \
	        fi; \
	    fi; \
	done;
	if ! [ -e site/$(HUGO_SITE)/config.yml ]; then \
	    touch site/$(HUGO_SITE)/config.yml; \
	    echo 'baseURL:        https://$(HUGO_SITE)/' >> site/$(HUGO_SITE)/config.yml; \
	    echo 'metaDataFormat: yaml'                  >> site/$(HUGO_SITE)/config.yml; \
	fi;

.PHONY: hugo-themes
hugo-themes:
	git clone --depth 1 --recursive https://github.com/spf13/hugoThemes.git themes

.PHONY: hugo-site
hugo-site:
	docker run --rm \
	    -v '$(CWD)/site':/usr/share \
	    -w /usr/share \
	    kamilsk/hugo:latest \
	    hugo new site $(HUGO_SITE)

.PHONY: hugo-theme
hugo-theme:
	docker run --rm \
	    -v '$(CWD)/site/$(HUGO_SITE)':/usr/share/site \
	    kamilsk/hugo:latest \
	    hugo new theme $(THEME)

.PHONY: hugo-content
hugo-content:
	docker run --rm \
	    -v '$(CWD)/site/$(HUGO_SITE)':/usr/share/site \
	    kamilsk/hugo:latest \
	    hugo new $(CONTENT).md

.PHONY: hugo-mount
hugo-mount:
	docker run --rm -it \
	    -v '$(CWD)/site/$(HUGO_SITE)':/usr/share/site \
	    -p $(HUGO_HOST):$(HUGO_PORT) \
	    -e PORT=$(HUGO_PORT) \
	    -e BASE_URL='http://$(HUGO_HOST)' \
	    kamilsk/hugo:latest /bin/sh

.PHONY: hugo-start
hugo-start:
	docker run --rm -d \
	    --name hugo-$(HUGO_SITE) \
	    -v '$(CWD)/site/$(HUGO_SITE)':/usr/share/site \
	    -p $(HUGO_HOST):$(HUGO_PORT) \
	    -e PORT=$(HUGO_PORT) \
	    -e BASE_URL='http://$(HUGO_HOST)' \
	    -e ARGS='$(strip $(ARGS))' \
	    kamilsk/hugo:latest

.PHONY: hugo-stop
hugo-stop:
	docker stop hugo-$(HUGO_SITE)

.PHONY: hugo-logs
hugo-logs:
	docker logs -f hugo-$(HUGO_SITE)
