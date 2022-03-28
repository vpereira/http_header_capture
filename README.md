### HTTP Header capture tool with data exfiltration via DNS

A weaponized version was used in a red-team exercise. We gained access to a server in a DMZ, which had access to SSO tokens used in different applications.

### Static building approach

We decided to go with a "static building" approach, linking statically libpcap. That avoid installing dependencies on the target server as much as dependencies incompatibilities. 
