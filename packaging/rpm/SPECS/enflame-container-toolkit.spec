Name: enflame-container-toolkit
Version: %{version}
Release: %{release}
Group: Development Tools

Vendor: NVIDIA CORPORATION
Packager: NVIDIA CORPORATION <topstools@enflame.com>

Summary: ENFLAME Container Toolkit
URL: https://enflame/enflame-container-toolkit
License: Apache-2.0

Source0: enflame-container-runtime-hook
Source1: enflame-ctk
Source2: LICENSE
Source3: enflame-container-runtime
Source4: enflame-container-runtime.cdi
Source5: enflame-container-runtime.legacy

Obsoletes: enflame-container-runtime <= 3.5.0-1, enflame-container-runtime-hook <= 1.4.0-2
Provides: enflame-container-runtime
Provides: enflame-container-runtime-hook
Requires: libenflame-container-tools >= %{libenflame_container_tools_version}, libenflame-container-tools < 2.0.0
Requires: enflame-container-toolkit-base == %{version}-%{release}

%description
Provides tools and utilities to enable GCU support in containers.

%prep
cp %{SOURCE0} %{SOURCE1} %{SOURCE2} %{SOURCE3} %{SOURCE4} %{SOURCE5} .

%install
mkdir -p %{buildroot}%{_bindir}
install -m 755 -t %{buildroot}%{_bindir} enflame-container-runtime-hook
install -m 755 -t %{buildroot}%{_bindir} enflame-container-runtime
install -m 755 -t %{buildroot}%{_bindir} enflame-container-runtime.cdi
install -m 755 -t %{buildroot}%{_bindir} enflame-container-runtime.legacy
install -m 755 -t %{buildroot}%{_bindir} enflame-ctk

%post
if [ $1 -gt 1 ]; then  # only on package upgrade
  mkdir -p %{_localstatedir}/lib/rpm-state/enflame-container-toolkit
  cp -af %{_bindir}/enflame-container-runtime-hook %{_localstatedir}/lib/rpm-state/enflame-container-toolkit
fi

%posttrans
if [ ! -e %{_bindir}/enflame-container-runtime-hook ]; then
  # repairing lost file enflame-container-runtime-hook
  cp -avf %{_localstatedir}/lib/rpm-state/enflame-container-toolkit/enflame-container-runtime-hook %{_bindir}
fi
rm -rf %{_localstatedir}/lib/rpm-state/enflame-container-toolkit
ln -sf %{_bindir}/enflame-container-runtime-hook %{_bindir}/enflame-container-toolkit

# Generate the default config; If this file already exists no changes are made.
%{_bindir}/enflame-ctk --quiet config --config-file=%{_sysconfdir}/enflame-container-runtime/config.toml --in-place

%postun
if [ "$1" = 0 ]; then  # package is uninstalled, not upgraded
  if [ -L %{_bindir}/enflame-container-toolkit ]; then rm -f %{_bindir}/enflame-container-toolkit; fi
fi

%files
%license LICENSE
%{_bindir}/enflame-container-runtime-hook

%changelog
# As of 1.10.0-1 we generate the release information automatically
* %{release_date} NVIDIA CORPORATION <topstools@enflame.com> %{version}-%{release}
- See https://gitlab.com/enflame/container-toolkit/container-toolkit/-/blob/%{git_commit}/CHANGELOG.md
- Bump libenflame-container dependency to libenflame-container-tools >= %{libenflame_container_tools_version}

# The BASE package consists of the ENFLAME Container Runtime and the ENFLAME Container Toolkit CLI.
# This allows the package to be installed on systems where no ENFLAME Container CLI is available.
%package base
Summary: ENFLAME Container Toolkit Base
Obsoletes: enflame-container-runtime <= 3.5.0-1, enflame-container-runtime-hook <= 1.4.0-2
Provides: enflame-container-runtime
# Since this package allows certain components of the ENFLAME Container Toolkit to be installed separately
# it conflicts with older versions of the enflame-container-toolkit package that also provide these files.
Conflicts: enflame-container-toolkit <= 1.10.0-1

%description base
Provides tools such as the ENFLAME Container Runtime and ENFLAME Container Toolkit CLI to enable GCU support in containers.

%files base
%license LICENSE
%{_bindir}/enflame-container-runtime
%{_bindir}/enflame-ctk

# The OPERATOR EXTENSIONS package consists of components that are required to enable GCU support in Kubernetes.
# This package is not distributed as part of the ENFLAME Container Toolkit RPMs.
%package operator-extensions
Summary: ENFLAME Container Toolkit Operator Extensions
Requires: enflame-container-toolkit-base == %{version}-%{release}

%description operator-extensions
Provides tools for using the ENFLAME Container Toolkit with the GCU Operator

%files operator-extensions
%license LICENSE
%{_bindir}/enflame-container-runtime.cdi
%{_bindir}/enflame-container-runtime.legacy
