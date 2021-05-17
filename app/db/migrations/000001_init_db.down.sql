BEGIN;
DROP TABLE public.accounts;
DROP TABLE public.configs;
DROP TABLE public.ops;
DROP TABLE public.projects;
DROP TABLE public.assets;
DROP TABLE public."project-collaborators";
DROP TABLE public.snapshots;

DROP TYPE enum_assets_type;
DROP TYPE enum_project-collaborators_accessType;
DROP TYPE enum_projects_accessType;
COMMIT;