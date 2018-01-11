package sub

func Reverse(s string) string {
	str :=[]rune(s)
	for i,j:=0,len(str)-1;i<len(str)/2;i,j=i+1,j-1{
		str[i],str[j]=str[j],str[i]
	}
	return string(str);
}