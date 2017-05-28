.PHONY: hugo-init
hugo-init:
	mkdir -p site/$(SITE)/archetypes \
	         site/$(SITE)/content \
	         site/$(SITE)/data \
	         site/$(SITE)/i18n \
	         site/$(SITE)/layouts \
	         site/$(SITE)/static \
	         site/$(SITE)/themes
	for file in $$(ls site/$(SITE)); do \
	    if [[ -d site/$(SITE)/$$file ]]; then \
	        export COUNT=$$(ls -a site/$(SITE)/$$file | wc -l); \
	        if [[ $$COUNT -lt 3 ]]; then \
	            touch site/$(SITE)/$$file/.gitkeep; \
	        fi; \
	    fi; \
	done;
	if ! [ -e site/$(SITE)/config.yml ]; then \
	    touch site/$(SITE)/config.yml; \
	    echo 'baseURL:        https://$(SITE)/' >> site/$(SITE)/config.yml; \
	    echo 'metaDataFormat: yaml'             >> site/$(SITE)/config.yml; \
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
	    hugo new site $(SITE)

.PHONY: hugo-theme
hugo-theme:
	docker run --rm \
	    -v '$(CWD)/site/$(SITE)':/usr/share/site \
	    kamilsk/hugo:latest \
	    hugo new theme $(THEME)

.PHONY: hugo-content
hugo-content:
	docker run --rm \
	    -v '$(CWD)/site/$(SITE)':/usr/share/site \
	    kamilsk/hugo:latest \
	    hugo new $(CONTENT).md

.PHONY: hugo-mount
hugo-mount:
	docker run --rm -it \
	    -v '$(CWD)/site/$(SITE)':/usr/share/site \
	    -p $(HOST):1313 \
	    kamilsk/hugo:latest /bin/sh

.PHONY: hugo-start
hugo-start:
	docker run --rm -d \
	    --name hugo-$(SITE) \
	    -v '$(CWD)/site/$(SITE)':/usr/share/site \
	    -p $(HOST):1313 \
	    -e ARGS='$(strip $(ARGS))' \
	    kamilsk/hugo:latest

.PHONY: hugo-stop
hugo-stop:
	docker stop hugo-$(SITE)

.PHONY: hugo-logs
hugo-logs:
	docker logs -f hugo-$(SITE)
