export const API_BASE_URL = "/api/v1";

export interface Link {
  title: string;
  url: string;
  type: string;
}

export interface Project {
  title: string;
  projectURL: string;
  projectImageURL: string;
  date: string;
}

export interface Affiliation {
  title: string;
  logoURL: string;
}

export interface Experience {
  current: boolean;
  end: string;
  start: string;
  company: string;
  position: string;
  logoURL: string;
  pageId: string;
  hightlights: string | undefined;
}

const monthNames = [
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];

function safeFetch<T>(
  url: string,
  options?: RequestInit
): Promise<[undefined, T] | [Error]> {
  return fetch(url, options)
    .then((response) => {
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      return response.json() as Promise<T>;
    })
    .then((data) => [undefined, data] as [undefined, T])
    .catch((error) => [error] as [Error]);
}

export const fetchLinks = async () => {
  const [error, response] = await safeFetch<{ links: Link[] }>(
    `${API_BASE_URL}/links`
  );
  if (error) {
    console.error("Error fetching links:", error);
    return { links: [] };
  }

  return response;
};

export const fetchAffiliations = async () => {
  const [error, response] = await safeFetch<{ affiliations: Affiliation[] }>(
    `${API_BASE_URL}/affiliations`
  );
  if (error) {
    console.error("Error fetching affiliations:", error);
    return { affiliations: [] };
  }

  return response;
};


export const fetchProjects = async () => {
  const [error, response] = await safeFetch<{ projects: Project[] }>(
    `${API_BASE_URL}/projects`
  );
  if (error) {
    console.error("Error fetching projects:", error);
    return { projects: [] };
  }

  return {
    projects: response.projects.map((project) => {
      const projectDate = new Date(project.date);
      return {
        ...project,
        date: `${monthNames[projectDate.getMonth()]} ${projectDate.getFullYear()}`,
      };
    }),
  };
};

export const fetchExperiences = async () => {
  const [error, response] = await safeFetch<{ experience: Experience[] }>(
    `${API_BASE_URL}/experience`
  );
  if (error) {
    console.error("Error fetching experiences:", error);
    return { experience: [] };
  }

  return {
    experience: response.experience.map((exp) => {
      const startDate = new Date(exp.start);
      const endDate = exp.current
        ? "Present"
        : new Date(exp.end).toLocaleDateString("en-US", {
            year: "numeric",
            month: "long",
          });
      return {
        ...exp,
        start: `${monthNames[startDate.getMonth()]} ${startDate.getFullYear()}`,
        end: endDate,
      };
    }),
  };
};

export const fetchContent = async (pageId: string) => {
  const [error, response] = await safeFetch<{ content: string }>(
    `${API_BASE_URL}/content/${pageId}`
  );
  if (error) {
    console.error("Error fetching content:", error);
    return { content: "" };
  }

  return {
    content: response.content,
  };
};