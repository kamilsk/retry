.PHONY: hugo-init
hugo-init: SITE = www
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
hugo-site: SITE = www
hugo-site:
	docker run --rm \
	    -v '$(CWD)/site':/opt \
	    -w /opt \
	    kamilsk/hugo:latest \
	    hugo new site $(SITE)

.PHONY: hugo-theme
hugo-theme: SITE  = www
hugo-theme: THEME = new
hugo-theme:
	docker run --rm \
	    -v '$(CWD)/site/$(SITE)':/opt \
	    -w /opt \
	    kamilsk/hugo:latest \
	    hugo new theme $(THEME)

.PHONY: hugo-content
hugo-content: SITE    = www
hugo-content: CONTENT = item
hugo-content:
	docker run --rm \
	    -v '$(CWD)/site/$(SITE)':/opt \
	    -w /opt \
	    kamilsk/hugo:latest \
	    hugo new $(CONTENT).md

.PHONY: hugo-mount
hugo-mount: SITE = www
hugo-mount:
	docker run --rm -it \
	    -v '$(CWD)/site/$(SITE)':/opt \
	    -w /opt \
	    -p 127.0.0.1:8080:8080 \
	    kamilsk/hugo:latest \
	    /bin/sh

.PHONY: hugo-start
hugo-start: SITE = www
hugo-start:
	docker run --rm -d \
	    -v '$(CWD)/site/$(SITE)':/opt \
	    -w /opt \
	    -p 127.0.0.1:8080:8080 \
	    kamilsk/hugo:latest \
	    /bin/sh -c 'hugo server --config=config.yml --baseURL=http://localhost:8080 --bind="" --port=8080 --buildDrafts'

.PHONY: hugo-stop
hugo-stop:
	docker ps | grep "/bin/sh -c 'hugo" | awk '{print $$1}' | xargs docker stop

.PHONY: hugo-logs
hugo-logs:
	docker ps | grep "/bin/sh -c 'hugo" | awk '{print $$1}' | xargs docker logs -f
