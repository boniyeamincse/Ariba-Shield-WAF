# Ariba Shield WAF — OpenResty gateway (Release 0.1)
FROM openresty/openresty:1.27.1.1-3-alpine

# Gateway runtime layout (ADR-003 write-ahead config store)
ENV CONFIG_STORE=/var/lib/shield-waf/config
ENV ACTIVE_POINTER=$CONFIG_STORE/active
ENV LKG_POINTER=$CONFIG_STORE/last-known-good
ENV CONFIG_SOURCE=/etc/nginx/conf.d/shield.conf

RUN mkdir -p $CONFIG_STORE /etc/nginx/conf.d /var/log/shield-waf \
 && addgroup -S shield 2>/dev/null || true \
 && adduser -S -G shield shield 2>/dev/null || true

# Copy runtime scripts and templates
COPY nginx/ /etc/nginx/shield/
COPY scripts/ /usr/local/bin/

RUN chmod +x /usr/local/bin/*.sh \
 && chown -R shield:shield $CONFIG_STORE /var/log/shield-waf

# Health check validates the running config is loadable
HEALTHCHECK --interval=10s --timeout=5s --retries=5 \
  CMD nginx -t

# Entrypoint performs last-known-good restore, then starts nginx
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
