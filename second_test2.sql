select departments.dept_name as 'department name',
	   titles.title,
	   employees.first_name as 'first name', 
       employees.last_name as 'last name',
	   employees.hire_date as 'hire date',
       year(curdate()) - year(employees.hire_date) as 'how many years they have been working'
from dept_emp, titles, employees, departments
where month(employees.hire_date) = month(curdate()) and
	  employees.emp_no  = titles.emp_no and
      employees.emp_no  = dept_emp.emp_no and
      titles.to_date >= curdate() and
      dept_emp.to_date >= curdate() and
      dept_emp.dept_no = departments.dept_no
      ;
